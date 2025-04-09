package file

import (
	"os"

	"github.com/neticdk/go-stdlib/file"
)

// SafeOpenFile opens a file with the specified path, base directory, flags, and mode.
// It ensures the file operation is secure by validating the mode and path.
// Returns a file handle and any error encountered.
// Deprecated: SafeOpenFile is deprecated - use github.com/neticdk/go-stdlib/file.SafeOpenFile
func SafeOpenFile(root, path string, flag int, mode int64) (*os.File, error) {
	return file.SafeOpenFile(root, path, flag, mode)
}

// SafeOpen opens a file for read-only access in a secure manner.
// It uses SafeOpenFile with read-only flag and default permissions.
// Deprecated: SafeOpen is deprecated - use github.com/neticdk/go-stdlib/file.SafeOpen
func SafeOpen(root, path string) (*os.File, error) {
	return file.SafeOpen(root, path)
}

// SafeCreate creates or truncates a file with the specified mode.
// It uses SafeOpenFile with write-only, create and truncate flags.
// Deprecated: SafeCreate is deprecated - use github.com/neticdk/go-stdlib/file.SafeCreate
func SafeCreate(root, path string, mode int64) (*os.File, error) {
	return file.SafeCreate(root, path, mode)
}

// SafeReadFile reads the entire contents of a file securely.
// It ensures the file path is safe before reading.
// Returns the file contents and any error encountered.
// Deprecated: SafeReadFile is deprecated - use github.com/neticdk/go-stdlib/file.SafeReadFile
func SafeReadFile(root, path string) ([]byte, error) {
	return file.SafeReadFile(root, path)
}

// ValidMode checks if the given mode is a valid file mode
// Deprecated: ValidMode is deprecated - use github.com/neticdk/go-stdlib/file.ValidMode
func ValidMode(mode int64) (os.FileMode, error) {
	return file.ValidMode(mode)
}

// SafePath ensures the given path is safe to use within the specified root directory.
// It returns the cleaned absolute path and an error if the path is unsafe.
// The function performs the following checks and operations:
// 1. Validates that the path is not empty and does not contain null bytes.
// 2. Cleans and converts both the root and input paths to absolute paths.
// 3. Resolves any symlinks in the input path, even if the path does not exist.
// 4. Ensures the resolved path is within the root directory to prevent path traversal attacks.
// 5. Works on both Windows and Unix-like systems, handling platform-specific path separators and case sensitivity.
//
// Parameters:
// - root: The root directory within which the path must be contained.
// - path: The input path to be validated and resolved. It can be either an absolute or relative path.
// Deprecated: SafePath is deprecated - use github.com/neticdk/go-stdlib/file.SafePath
func SafePath(root, path string) (string, error) {
	return file.SafePath(root, path)
}
