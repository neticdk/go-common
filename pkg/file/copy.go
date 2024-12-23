package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyDirectory copies a directory from src to dest
func CopyDirectory(srcDir, dest string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: '%s', error: '%s'", srcDir, err.Error())
	}
	if !IsDir(dest) {
		if err := os.MkdirAll(dest, 0o750); err != nil {
			return fmt.Errorf("failed to create directory: '%s', error: '%s'", dest, err.Error())
		}
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info for '%s'", sourcePath)
		}

		if err := copyFileOrDir(sourcePath, destPath, fileInfo); err != nil {
			return err
		}
	}
	return nil
}

// Copy copies a file from src to dest
func Copy(srcFile, dstFile string) error {
	in, err := SafeOpen(filepath.Dir(srcFile), srcFile)
	if err != nil {
		return fmt.Errorf("failed to open file: '%s', error: '%s'", srcFile, err.Error())
	}
	defer func() {
		if cerr := in.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: '%s', error: '%s'\n", srcFile, cerr.Error())
		}
	}()

	stat, err := in.Stat()
	var mode int64 = 0o640
	if err == nil {
		mode = int64(stat.Mode().Perm())
	}
	out, err := SafeCreate(filepath.Dir(dstFile), dstFile, mode)
	if err != nil {
		return fmt.Errorf("failed to create file: '%s', error: '%s'", dstFile, err.Error())
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: '%s', error: '%s'\n", dstFile, cerr.Error())
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy file: '%s', error: '%s'", srcFile, err.Error())
	}

	return nil
}

func copySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return fmt.Errorf("failed to read symlink: '%s', error: '%s'", source, err.Error())
	}
	return os.Symlink(link, dest)
}
