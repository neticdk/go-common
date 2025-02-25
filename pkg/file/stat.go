package file

import (
	"errors"
	"os"
)

// Exists returns true if the given path exists
//
// It returns false and an error on any error, e.g. on insufficient permissions
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

// IsDir returns true if the given path is a directory
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.IsDir()
}

// IsRegular returns true if the given path is a regular file
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsRegular(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.Mode().IsRegular()
}

// IsSymlink returns true if the given path is a symlink
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsSymlink(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink
}

// IsNamedPipe returns true if the given path is a named pipe
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsNamedPipe(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.Mode()&os.ModeNamedPipe == os.ModeNamedPipe
}

// IsSocket returns true if the given path is a socket
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsSocket(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.Mode()&os.ModeSocket == os.ModeSocket
}

// IsDevice returns true if the given path is a device
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsDevice(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.Mode()&os.ModeDevice == os.ModeDevice
}

// IsFile returns true if the given path is a regular file, symlink, socket, or device
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsFile(path string) bool {
	exists, err := Exists(path)
	if err != nil {
		return false
	}

	return exists && !IsDir(path)
}
