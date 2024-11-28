package helm

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// LoadChartFromFS loads a Helm chart from the filesystem
func LoadChartFromFS(filesystem fs.FS) (*chart.Chart, error) {
	var bufferedFiles []*loader.BufferedFile

	root := "."
	err := fs.WalkDir(filesystem, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk chart directory: %w", err)
		}
		if d.IsDir() {
			return nil
		}

		data, readErr := fs.ReadFile(filesystem, path)
		if readErr != nil {
			return fmt.Errorf("failed to read file %s: %w", path, readErr)
		}

		relativePath, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, relErr)
		}

		bufferedFile := &loader.BufferedFile{
			Name: relativePath,
			Data: data,
		}

		bufferedFiles = append(bufferedFiles, bufferedFile)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk chart directory: %w", err)
	}

	return loader.LoadFiles(bufferedFiles)
}
