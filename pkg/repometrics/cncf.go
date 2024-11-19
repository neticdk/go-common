package repometrics

import (
	"fmt"

	"github.com/neticdk/go-common/pkg/cncf"
)

// UpdateCNCFOptions represents the options for updating the CNCF status of a
//
// ProjectName, ProjectRepoURL, and ProjectHomepageURL are used to find the
// project.
//
// The order of precedence is:
// 1. ProjectRepoURL
// 2. ProjectHomepageURL
// 3. ProjectName (case-insensitive)
type UpdateCNCFOptions struct {
	Client             cncf.HTTPClient
	ProjectName        string
	ProjectRepoURL     string
	ProjectHomepageURL string
}

// UpdateCNCFStatus updates the CNCF status of the repository
func (m *Metrics) UpdateCNCFStatus(opts UpdateCNCFOptions) error {
	if opts.Client == nil {
		opts.Client = &cncf.DefaultHTTPClient{}
	}
	l, err := cncf.GetLandscape(opts.Client)
	if err != nil {
		return fmt.Errorf("failed to get landscape: %w", err)
	}

	findOpts := cncf.FindProjectOptions{
		RepoURL:     opts.ProjectRepoURL,
		HomepageURL: opts.ProjectHomepageURL,
		Name:        opts.ProjectName,
	}
	if project := l.FindProject(findOpts); project != nil {
		m.IsCNCF = true
		m.CNCFStatus = project.Project
	}

	return nil
}
