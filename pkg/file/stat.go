package file

import (
	"github.com/neticdk/go-stdlib/file"
)

// Exists returns true if the given path exists
//
// It returns false and an error on any error, e.g. on insufficient permissions
// Deprecated: Exists is deprecated - use github.com/neticdk/go-stdlib/file.Exists
func Exists(path string) (bool, error) {
	return file.Exists(path)
}

// IsDir returns true if the given path is a directory
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
// Deprecated: IsDir is deprecated - use github.com/neticdk/go-stdlib/file.IsDir
func IsDir(path string) bool {
	return file.IsDir(path)
}

// IsRegular returns true if the given path is a regular file
//
// It resolves all symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
// Deprecated: IsRegular is deprecated - use github.com/neticdk/go-stdlib/file.IsRegular
func IsRegular(path string) bool {
	return file.IsRegular(path)
}

// IsSymlink returns true if the given path is a symlink
//
// It does not resolve any symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
// Deprecated: IsSymlink is deprecated - use github.com/neticdk/go-stdlib/file.IsSymlink
func IsSymlink(path string) bool {
	return file.IsSymlink(path)
}

// IsNamedPipe returns true if the given path is a named pipe
//
// It resolves all symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
// Deprecated: IsNamedPipe is deprecated - use github.com/neticdk/go-stdlib/file.IsNamedPipe
func IsNamedPipe(path string) bool {
	return file.IsNamedPipe(path)
}

// IsSocket returns true if the given path is a socket
//
// It resolves all symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
// Deprecated: IsSocket is deprecated - use github.com/neticdk/go-stdlib/file.IsSocket
func IsSocket(path string) bool {
	return file.IsSocket(path)
}

// IsDevice returns true if the given path is a device
//
// It resolves all symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
// Deprecated: IsDevice is deprecated - use github.com/neticdk/go-stdlib/file.IsDevice
func IsDevice(path string) bool {
	return file.IsDevice(path)
}

// IsFile returns true if the given path is a regular file, symlink, socket, or device
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
// Deprecated: IsFile is deprecated - use github.com/neticdk/go-stdlib/file.IsFile
func IsFile(path string) bool {
	return file.IsFile(path)
}
