package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	gh "github.com/google/go-github/v82/github"
	ghmock "github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDownloader struct {
	mock.Mock
}

func (m *MockDownloader) Download(url, path string) (int64, error) {
	args := m.Called(url, path)
	return int64(args.Int(0)), args.Error(1)
}

func TestPullRelease(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		// defer os.RemoveAll(tmpDir)

		mockDownloader := new(MockDownloader)
		mockDownloader.On("Download", mock.Anything, mock.Anything).Return(0, nil)

		mockedHTTPClient := newMockedHTTPClient()

		art := &artifact.Artifact{
			Name:      "my-release",
			Version:   "v1.0.0",
			BaseDir:   tmpDir,
			AssetName: "my-asset.zip",
		}
		opts := &ReleaseOptions{
			Client:     gh.NewClient(mockedHTTPClient),
			Downloader: mockDownloader,
			Uncompress: func(src, dest string) error { return nil },
			Owner:      "example",
			Repository: "my-repo",
			AssetName:  "my-asset",
		}
		if err := os.MkdirAll(filepath.Join(tmpDir, art.Name), 0o750); err != nil {
			t.Fatal(err)
		}

		result, err := PullRelease(context.Background(), art, opts)
		require.NoError(t, err)
		assert.NotNil(t, result)

		artifactDir := filepath.Join(tmpDir, "my-asset")
		assert.DirExists(t, artifactDir)
	})

	t.Run("error_no_client", func(t *testing.T) {
		art := &artifact.Artifact{
			Name:    "my-release",
			Version: "v1.0.0",
		}
		opts := &ReleaseOptions{}

		_, err := PullRelease(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client is required")
	})

	t.Run("error_no_downloader", func(t *testing.T) {
		mockedHTTPClient := newMockedHTTPClient()

		art := &artifact.Artifact{
			Name:    "my-release",
			Version: "v1.0.0",
		}
		opts := &ReleaseOptions{
			Client: gh.NewClient(mockedHTTPClient),
		}

		_, err := PullRelease(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "downloader is required")
	})

	t.Run("error_get_release", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		mockDownloader := new(MockDownloader)
		mockDownloader.On("Download", mock.Anything, mock.Anything).Return(0, fmt.Errorf("download error"))

		mockedHTTPClient := ghmock.NewMockedHTTPClient()

		art := &artifact.Artifact{
			Name:    "my-release",
			Version: "v1.0.0",
			BaseDir: tmpDir,
		}
		opts := &ReleaseOptions{
			Client:     gh.NewClient(mockedHTTPClient),
			Downloader: mockDownloader,
			Owner:      "example",
			Repository: "my-repo",
			AssetName:  "my-asset",
		}

		_, err = PullRelease(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting release tagged")
	})

	t.Run("error_download_asset", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		mockDownloader := new(MockDownloader)
		mockDownloader.On("Download", mock.Anything, mock.Anything).Return(0, fmt.Errorf("download error"))

		mockedHTTPClient := newMockedHTTPClient()

		art := &artifact.Artifact{
			Name:    "my-release",
			Version: "v1.0.0",
			BaseDir: tmpDir,
		}
		opts := &ReleaseOptions{
			Client:     gh.NewClient(mockedHTTPClient),
			Downloader: mockDownloader,
			Owner:      "example",
			Repository: "my-repo",
			AssetName:  "my-asset",
		}

		_, err = PullRelease(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "downloading asset")
	})

	t.Run("error_uncompress_asset", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		mockDownloader := new(MockDownloader)
		mockDownloader.On("Download", mock.Anything, mock.Anything).Return(0, nil)

		mockedHTTPClient := newMockedHTTPClient()

		mockUncompress := func(src, dest string) error {
			return fmt.Errorf("uncompress error")
		}

		art := &artifact.Artifact{
			Name:    "my-release",
			Version: "v1.0.0",
			BaseDir: tmpDir,
		}
		opts := &ReleaseOptions{
			Client:     gh.NewClient(mockedHTTPClient),
			Downloader: mockDownloader,
			Uncompress: mockUncompress,
			Owner:      "example",
			Repository: "my-repo",
			AssetName:  "my-asset",
		}

		_, err = PullRelease(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "uncompressing asset")
	})
}

func newMockedHTTPClient() *http.Client {
	return ghmock.NewMockedHTTPClient(
		ghmock.WithRequestMatch(
			ghmock.GetReposReleasesTagsByOwnerByRepoByTag,
			gh.RepositoryRelease{
				Name: gh.Ptr("foobar"),
				Assets: []*gh.ReleaseAsset{
					{
						ID:                 gh.Ptr(int64(0)),
						Name:               gh.Ptr("my-asset.zip"),
						BrowserDownloadURL: gh.Ptr("https://example.com/my-asset"),
					},
				},
			},
		),
		ghmock.WithRequestMatch(
			ghmock.GetReposReleasesLatestByOwnerByRepo,
			gh.RepositoryRelease{
				Name: gh.Ptr("foobar"),
				Assets: []*gh.ReleaseAsset{
					{
						ID:                 gh.Ptr(int64(0)),
						Name:               gh.Ptr("my-asset.zip"),
						BrowserDownloadURL: gh.Ptr("https://example.com/my-asset"),
					},
				},
			},
		),
	)
}
