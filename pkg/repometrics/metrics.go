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
	Name                string     `json:"name"`                 // Name of the repository
	Type                RepoType   `json:"type"`                 // Type of the repository
	URL                 string     `json:"url"`                  // URL of the repository
	CreatedAt           *time.Time `json:"created_at"`           // Creation timestamp of the repository
	IsCNCF              bool       `json:"is_cncf"`              // True if the repository is part of CNCF
	CNCFStatus          string     `json:"cncf_status"`          // CNCF status of the repository
	IsKubernetesSIG     bool       `json:"is_kubernetes_sig"`    // True if the repository is a Kubernetes SIG project
	IsApache            bool       `json:"is_apache"`            // True if the repository is an Apache project
	BackingOrganization string     `json:"backing_organization"` // Organization backing the repository
	License             string     `json:"license"`              // License of the repository

	Version                        string     `json:"version"`            // The desired version to work with
	VersionCreatedAt               *time.Time `json:"version_created_at"` // Creation timestamp of the repository version
	VersionCriticalVulnerabilities int        `json:"criticals"`
	VersionHighVulnerabilities     int        `json:"highs"`
	VersionMediumVulnerabilities   int        `json:"mediums"`
	VersionLowVulnerabilities      int        `json:"lows"`

	FilteredVersionCriticalVulnerabilities int `json:"filteredCriticals"`
	FilteredVersionHighVulnerabilities     int `json:"filteredHighs"`
	FilteredVersionMediumVulnerabilities   int `json:"filteredMediums"`
	FilteredVersionLowVulnerabilities      int `json:"filteredLows"`

	FixedVersionCriticalVulnerabilities int `json:"fixedCriticals"`
	FixedVersionHighVulnerabilities     int `json:"fixedHighs"`
	FixedVersionMediumVulnerabilities   int `json:"fixedMediums"`
	FixedVersionLowVulnerabilities      int `json:"fixedLows"`

	LatestVersion                        string `json:"recentVersion"` // Current version of the repository
	LatestVersionCriticalVulnerabilities int    `json:"latestCriticals"`
	LatestVersionHighVulnerabilities     int    `json:"latestHighs"`
	LatestVersionMediumVulnerabilities   int    `json:"latestMediums"`
	LatestVersionLowVulnerabilities      int    `json:"latestLows"`

	Performance     string                   `json:"performance"`     // Performance of the repository
	Vulnerabilities []types.VulnerabilityCDX `json:"vulnerabilities"` // Vulnerabilities found in the repository
	Vex             []types.VEXCDX           `json:"vex"`             // Exploits found in the repository
	Stats           *Stats                   `json:"stats"`           // Statistics of the repository
	Conclusion      string                   `json:"conclusion"`      // Conclusion of the exploration
}

// New creates a new Metrics
func New(t RepoType) (*Metrics, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("invalid repository type: %s", t)
	}

	return &Metrics{
		Type:            t,
		Stats:           NewStats(),
		Vulnerabilities: make([]types.VulnerabilityCDX, 0),
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
