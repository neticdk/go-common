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

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.IsDir()
}

func IsFile(path string) bool {
	if exists, _ := Exists(path); !exists {
		return false
	}

	return !IsDir(path)
}
