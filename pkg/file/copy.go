package file

import (
	"github.com/neticdk/go-stdlib/file"
)

// CopyDirectory copies a directory from src to dest
// Deprecated: CopyDirectory is deprecated - use github.com/neticdk/go-stdlib/file.CopyDirectory
func CopyDirectory(srcDir, dest string) error {
	return file.CopyDirectory(srcDir, dest)
}

// Copy copies a file from src to dest
// Deprecated: Copy is deprecated - use github.com/neticdk/go-stdlib/file.Copy
func Copy(srcFile, dstFile string) error {
	return file.Copy(srcFile, dstFile)
}
