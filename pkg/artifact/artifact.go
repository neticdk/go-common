package artifact

import (
	"fmt"
	"os"
	"path/filepath"
)

// Artifact is a source artifact that can be pulled or pushed
type Artifact struct {
	// Name is the name of the artifact
	Name string

	// Repository is the repository of the source
	// e.g. the helm repository or the github repository
	Repository string

	// Branch is the branch of the source
	Branch string

	// Tag is the tag of the source
	Tag string

	// CommitHash is the commit hash of the source
	CommitHash string

	// Version is the version of the source
	// e.g. the version of the helm chart or github release
	Version string

	// URL is the URL of the source
	URL string

	// AssetName is the file name of the asset to be downloaded
	// typically used for archives
	AssetName string

	// SubDir is a subdirectory relative to the souce root directory that is
	// used for the basis in the component. Leading slashes are removed.
	// e.g. 'my/path' in the checkout of git URL
	// 'https://github.com/my/repo.git' would use 'my/path' as the directory
	// where the source is found
	SubDir string

	// BaseDir is the directory where the artifact is stored
	// It does not include the directory for the artifact itself
	BaseDir string
}

// DestDir returns the destination directory of the artifact
func (a *Artifact) DestDir() string {
	return filepath.Join(a.BaseDir, a.Name)
}

// DestDirExists checks if the destination directory exists
func (a *Artifact) DestDirExists() error {
	if _, err := os.Stat(a.DestDir()); os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("artifact directory already exists: %s", a.DestDir())
}

// FirstVersionOrLatest returns the first non-empty version from the list
// or "latest" if none are found
func FirstVersionOrLatest(versions ...string) string {
	for _, version := range versions {
		if version != "" {
			return version
		}
	}
	return "latest"
}
