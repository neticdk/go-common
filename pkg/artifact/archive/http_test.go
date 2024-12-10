package archive

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// a base64 encoded tar.gz file with a my-component-1.0.0/archive.txt file with
// the content "This is a test archive."
const goodArchiveContentB64 = "H4sIAL2kImcAA+2W207CMBzGi4kx4rVe9wWo/61dJxde4CmS4AmIBy4kDdZIdANHNfNNvNd38E18HougKM6hScUo/WVND2vX0772I/XgJtdoBe1WKEOVcwgQQIYBAN/z8HPMezG4rBf30hQ7zHfAc13XBwwO00mEjQ8kiauOEpEeynmYXk9XOz1Ned+bCn6N/wrTCzNoCqEt0cA7FXyI+3TL0KwOrg6XOnTzd1/7ZKFaLfeT3Ra3OswNVckMyuf130dEu30hSTtqXctQhA2JMlPo4f547XFdBgYmafmMXRFvSnEio8WfOwdG6t+BIf17HADh2NQA0phw/VPAgWoGctnxKVCXcgbE8Shz8zTvZD0fl4orhfLqZnF/ncRCqYgkyXW5sFcsKF6W5e3aRu0AsiyPK7pR6Sit0RuNZ397HSaVj6pfNN7HKP139fJe/1RnEPaMjySBCdd/wv6TuogaZ81rSVSsTPSh14Mz9h3/xylw6//GgvV/E02C/geW0NA5MFL/w/7PBaDM+r9xkOz/2JKnLTi1/u/fk6B/w7f/aP0DHfZ/HDi19/84qJ41O1g/AivZUfhl760eLRaL5Z/zBAR3BukAGgAA"

type mockDownloader struct {
	downloadData []byte
	downloadErr  error
}

func (m *mockDownloader) Download(url, path string) (int64, error) {
	if m.downloadErr != nil {
		return 0, m.downloadErr
	}
	err := os.WriteFile(path, m.downloadData, 0o640)
	if err != nil {
		return 0, err
	}
	return int64(len(m.downloadData)), nil
}

func TestPullHTTP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Create a mock archive
		archiveText := []byte("This is a test archive.\n")
		archiveData, _ := base64.StdEncoding.DecodeString(goodArchiveContentB64)
		mockDownloader := &mockDownloader{downloadData: archiveData}
		art := &artifact.Artifact{
			Name:    "my-component",
			Version: "1.0.0",
			URL:     "http://example.com/my-component-1.0.0.tar.gz",
			BaseDir: tmpDir,
		}
		opts := &HTTPOptions{Downloader: mockDownloader, Uncompress: Uncompress}

		result, err := PullHTTP(context.Background(), art, opts)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Verify the archive was uncompressed
		archiveDir := filepath.Join(tmpDir, "my-component")
		assert.DirExists(t, archiveDir)

		// Verify the archive contents
		archiveFile := filepath.Join(archiveDir, "archive.txt")
		data, err := os.ReadFile(archiveFile)
		require.NoError(t, err)
		assert.Equal(t, archiveText, data)
	})

	t.Run("error_downloader", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common--")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		mockDownloader := &mockDownloader{downloadErr: fmt.Errorf("download error")}
		art := &artifact.Artifact{
			Name:    "my-component",
			Version: "1.0.0",
			URL:     "http://example.com/my-component-1.0.0.tar.gz",
			BaseDir: tmpDir,
		}
		opts := &HTTPOptions{Downloader: mockDownloader, Uncompress: Uncompress}

		_, err = PullHTTP(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "download error")
	})

	t.Run("error_uncompress", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Create a mock archive
		archiveData, _ := base64.StdEncoding.DecodeString(goodArchiveContentB64)
		mockDownloader := &mockDownloader{downloadData: archiveData}

		// Mock uncompress function to return an error
		mockUncompress := func(src, dest string) error {
			return fmt.Errorf("uncompress error")
		}

		art := &artifact.Artifact{
			Name:    "my-component",
			Version: "1.0.0",
			URL:     "http://example.com/my-component-1.0.0.tar.gz",
			BaseDir: tmpDir,
		}
		opts := &HTTPOptions{Downloader: mockDownloader, Uncompress: mockUncompress}

		_, err = PullHTTP(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to uncompress asset")
	})
}

func Test_PullHTTP_NotCompressed(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go-common-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a mock archive
	archiveData := []byte("This is a test archive")
	mockDownloader := &mockDownloader{downloadData: archiveData}
	art := &artifact.Artifact{
		Name:    "my-component",
		Version: "1.0.0",
		URL:     "http://example.com/my-component-1.0.0.txt", // Not a compressed file
		BaseDir: tmpDir,
	}
	opts := &HTTPOptions{Downloader: mockDownloader, Uncompress: Uncompress}

	_, err = PullHTTP(context.Background(), art, opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a compressed file")
}

func TestPullHTTP_ExtractToExistingDir(t *testing.T) {
	t.Run("error_with_existing_file", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "go-common-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Create a mock archive
		archiveData, _ := base64.StdEncoding.DecodeString(goodArchiveContentB64)
		mockDownloader := &mockDownloader{downloadData: archiveData}
		art := &artifact.Artifact{
			Name:    "my-component",
			Version: "1.0.0",
			URL:     "http://example.com/my-component-1.0.0.tar.gz",
			BaseDir: tmpDir,
		}
		opts := &HTTPOptions{Downloader: mockDownloader, Uncompress: Uncompress}

		// Create an existing directory with a file
		existingDir := filepath.Join(tmpDir, "my-component")
		err = os.MkdirAll(existingDir, 0o750)
		require.NoError(t, err)
		existingFile := filepath.Join(existingDir, "existing.txt")
		err = os.WriteFile(existingFile, []byte("Existing file"), 0o640)
		require.NoError(t, err)

		result, err := PullHTTP(context.Background(), art, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to rename directory")
		assert.Nil(t, result)

		// Verify the existing file is still present
		data, err := os.ReadFile(existingFile)
		require.NoError(t, err)
		assert.Equal(t, []byte("Existing file"), data)
	})
}
