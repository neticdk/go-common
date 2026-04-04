package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateChecker_CheckForUpdateAsync(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/vnd.github.v3+json", r.Header.Get("Accept"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tag_name": "v1.2.0"}`))
	}))
	defer ts.Close()

	tempDir := t.TempDir()

	u := NewUpdateChecker(&ExecutionContext{AppName: "testapp", Version: "v1.0.0"}, "owner", "repo", "Run 'brew upgrade testapp'", WithCacheDuration(1*time.Hour))
	u.cacheDir = tempDir
	u.githubURL = ts.URL

	notifyChan := u.CheckForUpdateAsync()
	msg := <-notifyChan

	assert.Contains(t, msg, "v1.0.0 -> v1.2.0")
	assert.Contains(t, msg, "Run 'brew upgrade testapp'")

	// Check if cache is created
	cachePath := filepath.Join(tempDir, "update-owner-repo.json")
	assert.FileExists(t, cachePath)

	cacheData, err := os.ReadFile(cachePath)
	assert.NoError(t, err)

	var cache UpdateCache
	err = json.Unmarshal(cacheData, &cache)
	assert.NoError(t, err)
	assert.Equal(t, "v1.2.0", cache.LatestVersion)

	// Now run again to test cache hit
	u2 := NewUpdateChecker(&ExecutionContext{AppName: "testapp", Version: "1.1.0"}, "owner", "repo", "Run 'brew upgrade testapp'", WithCacheDuration(1*time.Hour))
	u2.cacheDir = tempDir
	u2.githubURL = "http://invalid.url" // Should not be hit

	notifyChan2 := u2.CheckForUpdateAsync()
	msg2 := <-notifyChan2

	assert.Contains(t, msg2, "1.1.0 -> v1.2.0")
}

func TestUpdateChecker_NoUpdate(t *testing.T) {
	tempDir := t.TempDir()

	cachePath := filepath.Join(tempDir, "update-owner-repo.json")
	_ = os.MkdirAll(tempDir, 0o755)
	newCache, _ := json.Marshal(UpdateCache{
		LastCheck:     time.Now(),
		LatestVersion: "v1.0.0",
	})
	_ = os.WriteFile(cachePath, newCache, 0o644)

	u := NewUpdateChecker(&ExecutionContext{AppName: "testapp", Version: "v1.0.0"}, "owner", "repo", "", WithCacheDuration(1*time.Hour))
	u.cacheDir = tempDir
	u.githubURL = "http://invalid.url"

	notifyChan := u.CheckForUpdateAsync()
	msg := <-notifyChan

	assert.Empty(t, msg)
}

func TestUpdateChecker_WithoutCache(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tag_name": "v1.3.0"}`))
	}))
	defer ts.Close()

	tempDir := t.TempDir()

	u := NewUpdateChecker(&ExecutionContext{AppName: "testapp", Version: "v1.0.0"}, "owner", "repo", "", WithoutCache())
	u.githubURL = ts.URL
	u.cacheDir = tempDir

	notifyChan := u.CheckForUpdateAsync()
	msg := <-notifyChan

	assert.Contains(t, msg, "v1.0.0 -> v1.3.0")

	// Verify that no cache file was created
	cachePath := filepath.Join(tempDir, "update-owner-repo.json")
	assert.NoFileExists(t, cachePath)
}

func TestUpdateChecker_WithMessageFormatter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tag_name": "v2.0.0"}`))
	}))
	defer ts.Close()

	u := NewUpdateChecker(&ExecutionContext{AppName: "testapp", Version: "v1.0.0"}, "owner", "repo", "", WithoutCache(), WithMessageFormatter(func(current, latest string) string {
		return fmt.Sprintf("Hey! %s is out (you have %s)", latest, current)
	}))
	u.githubURL = ts.URL

	notifyChan := u.CheckForUpdateAsync()
	msg := <-notifyChan

	assert.Equal(t, "Hey! v2.0.0 is out (you have v1.0.0)", msg)
}

func TestUpdateChecker_NoUpdateNotifierEnv(t *testing.T) {
	t.Setenv("NO_UPDATE_NOTIFIER", "1")

	u := NewUpdateChecker(&ExecutionContext{AppName: "testapp", Version: "v1.0.0"}, "owner", "repo", "")

	notifyChan := u.CheckForUpdateAsync()
	msg := <-notifyChan

	assert.Empty(t, msg)
}

func TestUpdateChecker_InvalidCurrentVersion(t *testing.T) {
	u := NewUpdateChecker(&ExecutionContext{AppName: "testapp", Version: "HEAD"}, "owner", "repo", "")

	notifyChan := u.CheckForUpdateAsync()
	msg := <-notifyChan

	assert.Empty(t, msg)
}

func TestEnsureVPrefix(t *testing.T) {
	assert.Equal(t, "v1.0.0", ensureVPrefix("1.0.0"))
	assert.Equal(t, "v1.0.0", ensureVPrefix("v1.0.0"))
	assert.Equal(t, "", ensureVPrefix(""))
}
