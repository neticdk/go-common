package file

import (
	"errors"
	"os"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// IsDir returns true if a directory exists at the path
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.IsDir()
}

// IsRegular returns true if a regular file exists at the path
func IsRegular(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.Mode().IsRegular()
}

// IsFile returns true if any type of file exists at the path
func IsFile(path string) bool {
	exists, err := Exists(path)
	if err != nil {
		return false
	}

	return exists && !IsDir(path)
}
