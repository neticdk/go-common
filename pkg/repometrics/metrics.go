package repometrics

import (
	"fmt"
	"time"

	"github.com/neticdk/go-common/pkg/types"
)

// perPage is the number of items to fetch per page when using pagination
const perPage = 1000

// Metrics represents the metrics of a repository
type Metrics struct {
	Name                string                `json:"name"`                 // Name of the repository
	Type                RepoType              `json:"type"`                 // Type of the repository
	URL                 string                `json:"url"`                  // URL of the repository
	CreatedAt           *time.Time            `json:"created_at"`           // Creation timestamp of the repository
	IsCNCF              bool                  `json:"is_cncf"`              // True if the repository is part of CNCF
	CNCFStatus          string                `json:"cncf_status"`          // CNCF status of the repository
	IsKubernetesSIG     bool                  `json:"is_kubernetes_sig"`    // True if the repository is a Kubernetes SIG project
	BackingOrganization string                `json:"backing_organization"` // Organization backing the repository
	License             string                `json:"license"`              // License of the repository
	CurrentVersion      string                `json:"current_version"`      // Current version of the repository
	Vulnerabilities     []types.Vulnerability `json:"vulnerabilities"`      // Vulnerabilities found in the repository
	Stats               *Stats                `json:"stats"`                // Statistics of the repository
}

// New creates a new Metrics
func New(t RepoType) (*Metrics, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("invalid repository type: %s", t)
	}

	return &Metrics{
		Type:            t,
		Stats:           NewStats(),
		Vulnerabilities: make([]types.Vulnerability, 0),
	}, nil
}

// ProjectAge returns the age of the project
func (m *Metrics) ProjectAge() time.Duration {
	if m.CreatedAt == nil {
		return 0
	}
	return time.Since(*m.CreatedAt)
}

type RepoType string

const (
	RepoTypeGitHub RepoType = "github"
)

func (t RepoType) IsValid() bool {
	switch t {
	case RepoTypeGitHub:
		return true
	default:
		return false
	}
}
