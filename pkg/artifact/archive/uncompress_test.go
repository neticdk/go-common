package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createZipArchive(t *testing.T, files map[string]string) []byte {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for name, content := range files {
		f, err := zipWriter.Create(name)
		require.NoError(t, err)
		_, err = f.Write([]byte(content))
		require.NoError(t, err)
	}

	err := zipWriter.Close()
	require.NoError(t, err)

	return buf.Bytes()
}

func createTarGzArchive(t *testing.T, files map[string]string) []byte {
	buf := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(gzipWriter)

	for name, content := range files {
		header := &tar.Header{
			Name: name,
			Mode: 0o640,
			Size: int64(len(content)),
		}
		err := tarWriter.WriteHeader(header)
		require.NoError(t, err)
		_, err = tarWriter.Write([]byte(content))
		require.NoError(t, err)
	}

	err := tarWriter.Close()
	require.NoError(t, err)
	err = gzipWriter.Close()
	require.NoError(t, err)

	return buf.Bytes()
}

func TestUncompress(t *testing.T) {
	t.Parallel()

	t.Run("unzip", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "uncompress-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		files := map[string]string{
			"file1.txt": "This is file 1",
			"file2.txt": "This is file 2",
		}
		zipData := createZipArchive(t, files)

		zipFilePath := filepath.Join(tmpDir, "test.zip")
		err = os.WriteFile(zipFilePath, zipData, 0o640)
		require.NoError(t, err)

		err = Uncompress(zipFilePath, tmpDir)
		require.NoError(t, err)

		for name, content := range files {
			filePath := filepath.Join(tmpDir, name)
			data, err := os.ReadFile(filePath)
			require.NoError(t, err)
			assert.Equal(t, content, string(data))
		}
	})

	t.Run("untarGz", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "uncompress-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		files := map[string]string{
			"file1.txt": "This is file 1",
			"file2.txt": "This is file 2",
		}
		tarGzData := createTarGzArchive(t, files)

		tarGzFilePath := filepath.Join(tmpDir, "test.tar.gz")
		err = os.WriteFile(tarGzFilePath, tarGzData, 0o640)
		require.NoError(t, err)

		err = Uncompress(tarGzFilePath, tmpDir)
		require.NoError(t, err)

		for name, content := range files {
			filePath := filepath.Join(tmpDir, name)
			data, err := os.ReadFile(filePath)
			require.NoError(t, err)
			assert.Equal(t, content, string(data))
		}
	})

	t.Run("unsupported_extension", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "uncompress-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		filePath := filepath.Join(tmpDir, "test.txt")
		err = os.WriteFile(filePath, []byte("This is a test file"), 0o640)
		require.NoError(t, err)

		err = Uncompress(filePath, tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported file extension")
	})
}

func TestIsArchive(t *testing.T) {
	assert.True(t, IsArchive("file.zip"))
	assert.True(t, IsArchive("file.tgz"))
	assert.True(t, IsArchive("file.tar.gz"))
	assert.False(t, IsArchive("file.txt"))
	assert.False(t, IsArchive("file.tar"))
}
