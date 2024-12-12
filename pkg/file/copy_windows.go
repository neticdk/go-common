//go:build windows

package file

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func copyFileOrDir(sourcePath, destPath string, fileInfo os.FileInfo) error {
	_, ok := fileInfo.Sys().(*windows.Win32FileAttributeData)
	if !ok {
		return fmt.Errorf("failed to get raw windows.Win32FileAttributeData data for '%s'", sourcePath)
	}

	switch fileInfo.Mode() & os.ModeType {
	case os.ModeDir:
		if err := os.MkdirAll(destPath, 0o750); err != nil {
			return fmt.Errorf("failed to create directory: '%s', error: '%s'", destPath, err.Error())
		}
		if err := CopyDirectory(sourcePath, destPath); err != nil {
			return fmt.Errorf("failed to copy directory: '%s', error: '%s'", sourcePath, err.Error())
		}
	case os.ModeSymlink:
		if err := copySymLink(sourcePath, destPath); err != nil {
			return fmt.Errorf("failed to copy symlink: '%s', error: '%s'", sourcePath, err.Error())
		}
	default:
		if err := Copy(sourcePath, destPath); err != nil {
			return fmt.Errorf("failed to copy file: '%s', error: '%s'", sourcePath, err.Error())
		}
	}

	isSymlink := fileInfo.Mode()&os.ModeSymlink != 0
	if !isSymlink {
		if err := os.Chmod(destPath, fileInfo.Mode()); err != nil {
			return fmt.Errorf("failed to change mode of '%s', error: '%s'", destPath, err.Error())
		}
	}

	// Not handling ownership and permissions on Windows yet

	return nil
}
