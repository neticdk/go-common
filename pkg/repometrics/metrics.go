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
	Name                string                `json:"name"`
	Type                RepoType              `json:"type"`
	URL                 string                `json:"url"`
	CreatedAt           *time.Time            `json:"created_at"`
	IsCNCF              bool                  `json:"is_cncf"`
	CNCFStatus          string                `json:"cncf_status"`
	IsKubernetesSIG     bool                  `json:"is_kubernetes_sig"`
	BackingOrganization string                `json:"backing_organization"`
	License             string                `json:"license"`
	CurrentVersion      string                `json:"current_version"`
	Vulnerabilities     []types.Vulnerability `json:"vulnerabilities"`
	Stats               *Stats                `json:"stats"`
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
