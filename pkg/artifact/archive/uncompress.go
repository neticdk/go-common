package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/neticdk/go-stdlib/file"
)

const (
	maxDecompressedSize = 1 * 1024 * 1024 * 1024 // 1GB
)

// Uncompress takes a file path and uncompresses it based on its extension.
// Supported extensions: zip, tgz, tar.gz
func Uncompress(filePath, dir string) error {
	f, err := file.SafeOpen(filepath.Dir(filePath), filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	switch {
	case strings.HasSuffix(filePath, ".zip"):
		s, err := f.Stat()
		if err != nil {
			return fmt.Errorf("getting file info: %w", err)
		}
		return unzip(f, s.Size(), dir)
	case strings.HasSuffix(filePath, ".tgz"), strings.HasSuffix(filePath, ".tar.gz"):
		return untarGz(f, dir)
	default:
		return fmt.Errorf("unsupported file extension")
	}
}

func unzip(f *os.File, size int64, dir string) error {
	// Create a new zip reader
	zr, err := zip.NewReader(f, size)
	if err != nil {
		return fmt.Errorf("creating zip reader: %w", err)
	}

	// Closure to address file descriptors issue with all the deferred .Close()
	// methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("opening file in zip: %w", err)
		}
		defer rc.Close()

		baseFile := filepath.Base(f.Name)
		path := filepath.Join(dir, baseFile)
		if f.FileInfo().Mode().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return fmt.Errorf("creating directory: %w", err)
			}
		} else if f.FileInfo().Mode().IsRegular() {
			if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
				return fmt.Errorf("creating directory: %w", err)
			}
			outFile, err := file.SafeOpenFile(dir, baseFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, int64(f.Mode()))
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer func() {
				if err := outFile.Close(); err != nil {
					panic(err)
				}
			}()

			limitReader := io.LimitReader(rc, maxDecompressedSize) // Limit the reader to 100MB

			_, err = io.Copy(outFile, limitReader)
			if err != nil {
				return fmt.Errorf("copying file contents: %w", err)
			}

			if stat, err := outFile.Stat(); err != nil {
				if stat.Size() > maxDecompressedSize {
					return fmt.Errorf("file too large")
				}
			}
		}
		return nil
	}

	// Iterate over the files in the archive
	for _, f := range zr.File {
		if err := extractAndWriteFile(f); err != nil {
			return err
		}
	}

	return nil
}

func untarGz(r io.Reader, dir string) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("creating gzip reader: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	limitReader := io.LimitReader(tr, maxDecompressedSize) // Limit the reader to 100MB

	// Closure to address file descriptors issue with all the deferred .Close()
	// methods
	extractAndWriteFile := func(header *tar.Header) error {
		baseFile := filepath.Base(header.Name)
		path := filepath.Join(dir, baseFile)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, file.FileModeNewDirectory); err != nil {
				return fmt.Errorf("creating directory: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), file.FileModeNewDirectory); err != nil {
				return fmt.Errorf("creating directory: %w", err)
			}
			outFile, err := file.SafeOpenFile(dir, baseFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.Mode)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, limitReader)
			if err != nil {
				return fmt.Errorf("copying file contents: %w", err)
			}

			if stat, err := outFile.Stat(); err != nil {
				if stat.Size() > maxDecompressedSize {
					return fmt.Errorf("file too large")
				}
			}
		}
		return nil
	}

	// Iterate over the files in the archive
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar header: %w", err)
		}

		if err := extractAndWriteFile(header); err != nil {
			return err
		}
	}

	return nil
}

func IsArchive(filePath string) bool {
	return strings.HasSuffix(filePath, ".zip") || strings.HasSuffix(filePath, ".tgz") || strings.HasSuffix(filePath, ".tar.gz")
}
