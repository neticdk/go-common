package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/neticdk/go-common/pkg/cli/ui"
	"golang.org/x/mod/semver"
)

const defaultCacheDuration = 24 // hours

// updateCache represents the data we store locally.
type updateCache struct {
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
	cacheEnabled        bool
	disabled            bool
	releaseNameFormat   string
	messageFormatter    UpdateMessageFormatter
	logger              *slog.Logger
}

// UpdateCheckerOption is a function that configures an UpdateChecker
type UpdateCheckerOption func(*UpdateChecker)

// WithCacheDuration sets the duration for which the update check result is cached
func WithCacheDuration(d time.Duration) UpdateCheckerOption {
	return func(u *UpdateChecker) {
		u.cacheDuration = d
	}
}

// WithCache enables or disables caching for the update checker
func WithCache(enabled bool) UpdateCheckerOption {
	return func(u *UpdateChecker) {
		u.cacheEnabled = enabled
	}
}

// WithReleaseNameFormat sets a format string (e.g., "myapp-%s") to find releases by tag.
// When set, it fetches the list of releases and finds the most recent one matching the format,
// extracting the version in place of "%s".
func WithReleaseNameFormat(format string) UpdateCheckerOption {
	return func(u *UpdateChecker) {
		u.releaseNameFormat = format
		u.githubURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", url.PathEscape(u.githubOwner), url.PathEscape(u.githubRepo))
	}
}

// WithMessageFormatter allows setting a custom function to format the update message
func WithMessageFormatter(f UpdateMessageFormatter) UpdateCheckerOption {
	return func(u *UpdateChecker) {
		u.messageFormatter = f
	}
}

// UIMessageFormatter returns an UpdateMessageFormatter that formats the message using the ui package.
func UIMessageFormatter(installInstructions string) UpdateMessageFormatter {
	return func(current, latest string) string {
		msg := "\n" + ui.Info.Sprintf("A new version is available! (%s -> %s)", current, latest)
		if installInstructions != "" {
			msg += fmt.Sprintf("\n%s", installInstructions)
		}
		return msg
	}
}

// WithAppName overrides the default app name.
// This is useful when creating an update checker for secondary dependencies.
func WithAppName(name string) UpdateCheckerOption {
	return func(u *UpdateChecker) {
		u.appName = name
		dir, err := os.UserCacheDir()
		if err != nil {
			dir = os.TempDir()
		}
		u.cacheDir = filepath.Join(dir, name)
	}
}

// WithCurrentVersion overrides the default current version.
func WithCurrentVersion(version string) UpdateCheckerOption {
	return func(u *UpdateChecker) {
		u.currentVersion = version
	}
}

// NewUpdateChecker creates a new UpdateChecker
func NewUpdateChecker(ec *ExecutionContext, githubOwner, githubRepo, installInstructions string, opts ...UpdateCheckerOption) *UpdateChecker {
	appName := "unknown"
	version := ""
	var logger *slog.Logger
	if ec != nil {
		appName = ec.AppName
		version = ec.Version
		logger = ec.Logger
	}

	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	u := &UpdateChecker{
		appName:             appName,
		currentVersion:      version,
		githubOwner:         githubOwner,
		githubRepo:          githubRepo,
		githubURL:           fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", url.PathEscape(githubOwner), url.PathEscape(githubRepo)),
		installInstructions: installInstructions,
		cacheDir:            filepath.Join(dir, appName),
		cacheDuration:       defaultCacheDuration * time.Hour,
		cacheEnabled:        true,
		disabled:            os.Getenv("NO_UPDATE_NOTIFIER") != "",
		logger:              logger,
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

	if u.disabled {
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

		cachePath := filepath.Join(u.cacheDir, fmt.Sprintf("update-%s-%s.json", u.githubOwner, u.githubRepo))

		if u.cacheEnabled {
			cacheData, err := os.ReadFile(cachePath)
			var cache updateCache
			if err == nil && json.Unmarshal(cacheData, &cache) == nil && time.Since(cache.LastCheck) < u.cacheDuration {
				latestVersion = cache.LatestVersion
			}
		}

		if latestVersion == "" {
			var err error
			latestVersion, err = u.fetchLatestFromGitHub()
			if err != nil && u.logger != nil {
				u.logger.Debug("Update check failed", "error", err)
			}
			if latestVersion != "" && u.cacheEnabled {
				_ = os.MkdirAll(filepath.Dir(cachePath), 0o755)
				newCache, _ := json.Marshal(updateCache{
					LastCheck:     time.Now(),
					LatestVersion: latestVersion,
				})
				_ = os.WriteFile(cachePath, newCache, 0o640)
			}
		}

		if latestVersion == "" {
			if u.logger != nil {
				u.logger.Debug("Update check aborted: no latest version found")
			}
			return
		}

		latestSemver := ensureVPrefix(latestVersion)

		if semver.IsValid(latestSemver) && semver.Compare(currentSemver, latestSemver) < 0 {
			notifyChan <- u.messageFormatter(u.currentVersion, latestVersion)
		}
	}()

	return notifyChan
}

func (u *UpdateChecker) fetchLatestFromGitHub() (string, error) {
	client := http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequest(http.MethodGet, u.githubURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", u.appName)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	if u.releaseNameFormat != "" {
		var releases []struct {
			TagName    string `json:"tag_name"`
			Prerelease bool   `json:"prerelease"`
			Draft      bool   `json:"draft"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
			return "", err
		}

		parts := strings.Split(u.releaseNameFormat, "%s")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid release name format: %s", u.releaseNameFormat)
		}
		prefix, suffix := parts[0], parts[1]

		for _, r := range releases {
			if r.Draft || r.Prerelease {
				continue
			}
			if strings.HasPrefix(r.TagName, prefix) && strings.HasSuffix(r.TagName, suffix) {
				version := strings.TrimSuffix(strings.TrimPrefix(r.TagName, prefix), suffix)
				return version, nil
			}
		}
		return "", nil
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func ensureVPrefix(v string) string {
	if v != "" && v[0] != 'v' {
		return "v" + v
	}
	return v
}
