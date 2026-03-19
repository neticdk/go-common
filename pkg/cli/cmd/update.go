package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/mod/semver"
)

// UpdateCache represents the data we store locally.
type UpdateCache struct {
	LastCheck     time.Time `json:"last_check"`
	LatestVersion string    `json:"latest_version"`
}

// UpdateMessageFormatter is a function type for formatting the update message
type UpdateMessageFormatter func(currentVersion, latestVersion string) string

// UpdateChecker is used to check if a newer version of the CLI is available
type UpdateChecker struct {
	appName             string
	currentVersion      string
	githubOwner         string
	githubRepo          string
	githubURL           string
	installInstructions string
	cacheDir            string
	cacheDuration       time.Duration
	disableCache        bool
	messageFormatter    UpdateMessageFormatter
}

// UpdateCheckerOption is a function that configures an UpdateChecker
type UpdateCheckerOption func(*UpdateChecker)

// WithCacheDuration sets the duration for which the update check result is cached
func WithCacheDuration(d time.Duration) UpdateCheckerOption {
	return func(u *UpdateChecker) {
		u.cacheDuration = d
	}
}

// WithoutCache disables caching for the update checker
func WithoutCache() UpdateCheckerOption {
	return func(u *UpdateChecker) {
		u.disableCache = true
	}
}

// WithMessageFormatter allows setting a custom function to format the update message
func WithMessageFormatter(f UpdateMessageFormatter) UpdateCheckerOption {
	return func(u *UpdateChecker) {
		u.messageFormatter = f
	}
}

// NewUpdateChecker creates a new UpdateChecker
func NewUpdateChecker(ec *ExecutionContext, githubOwner, githubRepo, installInstructions string, opts ...UpdateCheckerOption) *UpdateChecker {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	u := &UpdateChecker{
		appName:             ec.AppName,
		currentVersion:      ec.Version,
		githubOwner:         githubOwner,
		githubRepo:          githubRepo,
		githubURL:           fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubOwner, githubRepo),
		installInstructions: installInstructions,
		cacheDir:            filepath.Join(dir, ec.AppName),
		cacheDuration:       24 * time.Hour,
		messageFormatter: func(current, latest string) string {
			msg := fmt.Sprintf("\n🚀 A new version is available! (%s -> %s)", current, latest)
			if installInstructions != "" {
				msg += fmt.Sprintf("\n%s", installInstructions)
			}
			return msg
		},
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

// CheckForUpdateAsync checks for an update asynchronously and returns a channel that will receive the update message if one is available
func (u *UpdateChecker) CheckForUpdateAsync() <-chan string {
	notifyChan := make(chan string, 1)

	if os.Getenv("NO_UPDATE_NOTIFIER") != "" {
		close(notifyChan)
		return notifyChan
	}

	currentSemver := ensureVPrefix(u.currentVersion)
	if !semver.IsValid(currentSemver) {
		close(notifyChan)
		return notifyChan
	}

	go func() {
		defer close(notifyChan)

		var latestVersion string

		cachePath := filepath.Join(u.cacheDir, "update.json")

		if !u.disableCache {
			cacheData, err := os.ReadFile(cachePath)
			var cache UpdateCache
			if err == nil && json.Unmarshal(cacheData, &cache) == nil && time.Since(cache.LastCheck) < u.cacheDuration {
				latestVersion = cache.LatestVersion
			}
		}

		if latestVersion == "" {
			latestVersion = u.fetchLatestFromGitHub()
			if latestVersion != "" && !u.disableCache {
				_ = os.MkdirAll(filepath.Dir(cachePath), 0o755)
				newCache, _ := json.Marshal(UpdateCache{
					LastCheck:     time.Now(),
					LatestVersion: latestVersion,
				})
				_ = os.WriteFile(cachePath, newCache, 0o644)
			}
		}

		if latestVersion == "" {
			return // Network failed or other error, silently exit
		}

		latestSemver := ensureVPrefix(latestVersion)

		if semver.IsValid(latestSemver) && semver.Compare(currentSemver, latestSemver) < 0 {
			if u.messageFormatter != nil {
				notifyChan <- u.messageFormatter(u.currentVersion, latestVersion)
			} else {
				msg := fmt.Sprintf("\n🚀 A new version is available! (%s -> %s)", u.currentVersion, latestVersion)
				if u.installInstructions != "" {
					msg += fmt.Sprintf("\n%s", u.installInstructions)
				}
				notifyChan <- msg
			}
		}
	}()

	return notifyChan
}

func (u *UpdateChecker) fetchLatestFromGitHub() string {
	client := http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequest(http.MethodGet, u.githubURL, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", u.appName)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return ""
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return ""
	}

	return release.TagName
}

func ensureVPrefix(v string) string {
	if v != "" && v[0] != 'v' {
		return "v" + v
	}
	return v
}
