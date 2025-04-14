package repometrics

import (
	"math"
	"sort"
	"strings"
	"time"

	vul "github.com/anchore/grype/grype/db/v6"
)

// Stats represents the metrics of a repository
// TODO: maybe separate these into individual structs.
type Stats struct {
	LastCommit                *time.Time `json,yaml:"last_commit"`                   // Date of the last commit.
	CommitsPerMonth6M         int        `json,yaml:"commits_per_month_6m"`          // Number of commits in the last 6 months.
	VerifiedCommitsPerMonth6M int        `json,yaml:"verified_commits_per_month_6m"` // Number of verified commits in the last 6 months.
	Contributors1Y            int        `json,yaml:"contributors_1y"`               // Number of contributors in the last year.
	FirstRelease              *time.Time `json,yaml:"first_release"`                 // Date of the first release.
	LastRelease               *time.Time `json,yaml:"last_release"`                  // Date of the last release.
	NoOfReleases              int        `json,yaml:"releases"`                      // Total number of releases.
	ReleasesPerDay            float64    `json,yaml:"releases_per_day"`              // Number of releases per day.
	ReleasesPerWeek           float64    `json,yaml:"releases_per_week"`             // Number of releases per week.
	ReleasesPerMonth          float64    `json,yaml:"releases_per_month"`            // Number of releases per month.
	ReleasesPerYear           float64    `json,yaml:"releases_per_year"`             // Number of releases per year.
	Releases                  []Release  `json,yaml:"releases"`                      // List of releases.

	OpenedIssuesNow   int `json,yaml:"open_issues_now"`     // Number of currently open issues.
	ClosedIssuesNow   int `json,yaml:"closed_issues_now"`   // Number of currently closed issues.
	OpenedPRsNow      int `json,yaml:"opened_prs_now"`      // Number of pull requests opened in the last month.
	ClosedPRsNow      int `json,yaml:"closed_prs_now"`      // Number of pull requests closed in the last month.
	OpenedFeaturesNow int `json,yaml:"opened_features_now"` // Number of issues opened in the last month.
	ClosedFeaturesNow int `json,yaml:"closed_features_now"` // Number of issues closed in the last month.
	OpenedBugsNow     int `json,yaml:"opened_bugs_now"`     // Number of pull requests opened in the last month.
	ClosedBugsNow     int `json,yaml:"closed_bugs_now"`     // Number of pull requests closed in the last month.

	OpenedIssues1M   int `json,yaml:"opened_issues_1m"`   // Number of issues opened in the last month.
	ClosedIssues1M   int `json,yaml:"closed_issues_1m"`   // Number of issues closed in the last month.
	OpenedPRs1M      int `json,yaml:"opened_prs_1m"`      // Number of pull requests opened in the last month.
	ClosedPRs1M      int `json,yaml:"closed_prs_1m"`      // Number of pull requests closed in the last month.
	OpenedFeatures1M int `json,yaml:"opened_features_1m"` // Number of issues opened in the last month.
	ClosedFeatures1M int `json,yaml:"closed_features_1m"` // Number of issues closed in the last month.
	OpenedBugs1M     int `json,yaml:"opened_bugs_1m"`     // Number of pull requests opened in the last month.
	ClosedBugs1M     int `json,yaml:"closed_bugs_1m"`     // Number of pull requests closed in the last month.

	OpenedIssues3M   int `json,yaml:"opened_issues_3m"`   // Number of issues opened in the last 3 months.
	ClosedIssues3M   int `json,yaml:"closed_issues_3m"`   // Number of issues closed in the last 3 months.
	OpenedPRs3M      int `json,yaml:"opened_prs_3m"`      // Number of pull requests opened in the last 3 months.
	ClosedPRs3M      int `json,yaml:"closed_prs_3m"`      // Number of pull requests closed in the last 3 months.
	OpenedFeatures3M int `json,yaml:"opened_features_3m"` // Number of issues opened in the last 3 months.
	ClosedFeatures3M int `json,yaml:"closed_features_3m"` // Number of issues closed in the last 3 months.
	OpenedBugs3M     int `json,yaml:"opened_bugs_3m"`     // Number of pull requests opened in the last 3 months.
	ClosedBugs3M     int `json,yaml:"closed_bugs_3m"`     // Number of pull requests closed in the last 3 months.

	OpenedIssues6M   int `json,yaml:"opened_issues_6m"`   // Number of issues opened in the last 6 months.
	ClosedIssues6M   int `json,yaml:"closed_issues_6m"`   // Number of issues closed in the last 6 months.
	OpenedPRs6M      int `json,yaml:"opened_prs_6m"`      // Number of pull requests opened in the last 6 months.
	ClosedPRs6M      int `json,yaml:"closed_prs_6m"`      // Number of pull requests closed in the last 6 months.
	OpenedFeatures6M int `json,yaml:"opened_features_6m"` // Number of issues opened in the last 6 months.
	ClosedFeatures6M int `json,yaml:"closed_features_6m"` // Number of issues closed in the last 6 months.
	OpenedBugs6M     int `json,yaml:"opened_bugs_6m"`     // Number of pull requests opened in the last 6 months.
	ClosedBugs6M     int `json,yaml:"closed_bugs_6m"`     // Number of pull requests closed in the last 6 months.

	OpenedIssues9M   int `json,yaml:"opened_issues_9m"`   // Number of issues opened in the last 9 months.
	ClosedIssues9M   int `json,yaml:"closed_issues_9m"`   // Number of issues closed in the last 9 months.
	OpenedPRs9M      int `json,yaml:"opened_prs_9m"`      // Number of pull requests opened in the last 9 months.
	ClosedPRs9M      int `json,yaml:"closed_prs_9m"`      // Number of pull requests closed in the last 9 months.
	OpenedFeatures9M int `json,yaml:"opened_features_9m"` // Number of issues opened in the last 9 months.
	ClosedFeatures9M int `json,yaml:"closed_features_9m"` // Number of issues closed in the last 9 months.
	OpenedBugs9M     int `json,yaml:"opened_bugs_9m"`     // Number of pull requests opened in the last 9 months.
	ClosedBugs9M     int `json,yaml:"closed_bugs_9m"`     // Number of pull requests closed in the last 9 months.

	OpenedIssues1Y   int `json,yaml:"opened_issues_1y"`   // Number of issues opened in the last year.
	ClosedIssues1Y   int `json,yaml:"closed_issues_1y"`   // Number of issues closed in the last year.
	OpenedPRs1Y      int `json,yaml:"opened_prs_1y"`      // Number of pull requests opened in the last year.
	ClosedPRs1Y      int `json,yaml:"closed_prs_1y"`      // Number of pull requests closed in the last year.
	OpenedFeatures1Y int `json,yaml:"opened_features_1y"` // Number of issues opened in the last year.
	ClosedFeatures1Y int `json,yaml:"closed_features_1y"` // Number of issues closed in the last year.
	OpenedBugs1Y     int `json,yaml:"opened_bugs_1y"`     // Number of pull requests opened in the last year.
	ClosedBugs1Y     int `json,yaml:"closed_bugs_1y"`     // Number of pull requests closed in the last year.

	Likes           int           `json,yaml:"likes"`             // Number of likes/hearts/stars.
	Forks           int           `json,yaml:"forks"`             // Number of forks.
	TopCommitters   []Contributor `json,yaml:"top_committers"`    // Top committers overall.
	TopCommitters1Y []Contributor `json,yaml:"top_committers_1y"` // Top committers in the last year.
	TopCommitters9M []Contributor `json,yaml:"top_committers_9m"` // Top committers in the last 9 months.
	TopCommitters6M []Contributor `json,yaml:"top_committers_6m"` // Top committers in the last 6 months.
	TopCommitters3M []Contributor `json,yaml:"top_committers_3m"` // Top committers in the last 3 months.
	TopCommitters1M []Contributor `json,yaml:"top_committers_1m"` // Top committers in the last 1 month.

	InactiveContributors1Y []Contributor `json,yaml:"inactive_contributors"` // Inactive contributors.
	InactiveContributors9M []Contributor `json,yaml:"inactive_contributors"` // Inactive contributors.
	InactiveContributors6M []Contributor `json,yaml:"inactive_contributors"` // Inactive contributors.
	InactiveContributors3M []Contributor `json,yaml:"inactive_contributors"` // Inactive contributors.

	VulnerabilitiesIndex1Y int `json,yaml:"vulnerabilities_index_1y"` // Vulnerabilities index in the last year.
	VulnerabilitiesIndex9M int `json,yaml:"vulnerabilities_index_9m"` // Vulnerabilities index in the last 9 months.
	VulnerabilitiesIndex6M int `json,yaml:"vulnerabilities_index_6m"` // Vulnerabilities index in the last 6 months.
	VulnerabilitiesIndex3M int `json,yaml:"vulnerabilities_index_3m"` // Vulnerabilities index in the last 3 months.
	VulnerabilitiesIndex1M int `json,yaml:"vulnerabilities_index_1m"` // Vulnerabilities index in the last 1 month.

	ReleasesIndex1Y int `json,yaml:"releases_index_1y"` // Releases index in the last year.
	ReleasesIndex9M int `json,yaml:"releases_index_9m"` // Releases index in the last 9 months.
	ReleasesIndex6M int `json,yaml:"releases_index_6m"` // Releases index in the last 6 months.
	ReleasesIndex3M int `json,yaml:"releases_index_3m"` // Releases index in the last 3 months.
	ReleasesIndex1M int `json,yaml:"releases_index_1m"` // Releases index in the last 1 month.
}

// Contributor represents a contributor to a repository
type Contributor struct {
	Name    string `json,yaml:"name"`
	URL     string `json,yaml:"url"`
	Commits int    `json,yaml:"commits"`
}

type Release struct {
	AssetsURL   string    `json:"assets_url,omitempty"`        // URL of the release assets.
	Date        time.Time `json,yaml:"date"`                   // Date of the release.
	Name        string    `json,yaml:"name"`                   // Name of the release.
	TarballURL  string    `json, yaml:"tarball_url,omitempty"` // URL of the release tarball.
	VulnRefs    ExtRefs   `json,yaml:"refs,omitempty""`        // List of references to SBOM, VEX and Vulnerabilities.
	MissingRefs ExtRefs   `json,yaml:"missed,omitempty""`      // List of missing references to SBOM, VEX and Vulnerabilities.
	ReleaseURL  string    `json,yaml:"url,omitempty""`         // URL of the release.
	UploadURL   string    `json,yaml:"uploaded,omitempty"`     // URL of the release upload.
}

type ExtRefs struct {
	SBOMCDX string `json,yaml:"sbomCDX"`
	VEXCDX  string `json,yaml:"vexCDX"`
	VULNCDX string `json,yaml:"vulnerabilityCDX"`
	VULNMD  string `json,yaml:"vulnerabilityMD"`
	Latest  bool   `json,yaml:"actual"`
}

type CVE struct {
	Vulnerability  vul.ID   `json,yaml:"id"`             // the Vulnerability
	Exploitability string   `json,yaml:"exploitability"` // the Exploitability not correct type yet
	SBOMs          []string `json,yaml:"sboms"`          // the SBOMs not correct type yet
}

func NewStats() *Stats {
	return &Stats{
		TopCommitters:   make([]Contributor, 0),
		TopCommitters1Y: make([]Contributor, 0),
	}
}

func sortContributors(contributors []Contributor) {
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Commits > contributors[j].Commits
	})
}

func getTopContributors(contributors []Contributor, limit int) []Contributor {
	if len(contributors) > limit {
		return contributors[:limit]
	}
	return contributors
}

func isBot(login string) bool {
	return strings.Contains(login, "[bot]") ||
		strings.Contains(login, "bot") ||
		strings.Contains(login, "istio-testing") ||
		strings.Contains(login, "fluxcdbot") ||
		strings.Contains(login, "dependabot") ||
		strings.Contains(login, "renovate-bot") ||
		strings.Contains(login, "renovate[bot]") ||
		strings.Contains(login, "action")
}

type issueCounts struct {
	openedIssues   int
	closedIssues   int
	openedPulls    int
	closedPulls    int
	openedFeatures int
	closedFeatures int
	openedBugs     int
	closedBugs     int
}

type ReleaseMetrics struct {
	PerDay   float64
	PerWeek  float64
	PerMonth float64
	PerYear  float64
}

const (
	daysInWeek  = 7
	daysInMonth = 30.44
	daysInYear  = 365.25
)

// CalculateReleaseMetrics calculates the release metrics based on the first and
// last release dates and the total number of releases.
func CalculateReleaseMetrics(firstRelease, lastRelease *time.Time, totalReleases int) ReleaseMetrics {
	if firstRelease == nil || lastRelease == nil || totalReleases == 0 {
		return ReleaseMetrics{}
	}

	period := lastRelease.Sub(*firstRelease)
	if period <= 0 {
		return ReleaseMetrics{}
	}

	daysInPeriod := period.Hours() / 24
	if daysInPeriod == 0 {
		return ReleaseMetrics{}
	}

	weeksInPeriod := daysInPeriod / daysInWeek
	monthsInPeriod := daysInPeriod / daysInMonth
	yearsInPeriod := daysInPeriod / daysInYear

	releasesPerDay := float64(totalReleases) / daysInPeriod
	releasesPerWeek := float64(totalReleases) / weeksInPeriod
	releasesPerMonth := float64(totalReleases) / monthsInPeriod
	releasesPerYear := float64(totalReleases) / yearsInPeriod

	return ReleaseMetrics{
		PerDay:   math.Round(releasesPerDay*10) / 10,
		PerWeek:  math.Round(releasesPerWeek*10) / 10,
		PerMonth: math.Round(releasesPerMonth*100) / 100,
		PerYear:  math.Round(releasesPerYear*100) / 100,
	}
}
