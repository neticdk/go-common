// Deprecated: version is deprecated and has been moved to github.com/neticdk/go-stdlib/version
package version

import "github.com/neticdk/go-stdlib/version"

// First returns the first non-empty version string from the provided list of
// versions.
// Deprecated: First is deprecated and has been moved to github.com/neticdk/go-stdlib/version.First
func First(versions ...string) string {
	return version.First(versions...)
}
