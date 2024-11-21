package repometrics

import (
	"math"
	"sort"
	"strings"
	"time"
)

// Stats represents the metrics of a repository
type Stats struct {
	LastCommit        *time.Time    `json:"last_commit"`          // Date of the last commit.
	CommitsPerMonth6M int           `json:"commits_per_month_6m"` // Number of commits in the last 6 months.
	Contributors1Y    int           `json:"contributors_1y"`      // Number of contributors in the last year.
	FirstRelease      *time.Time    `json:"first_release"`        // Date of the first release.
	LastRelease       *time.Time    `json:"last_release"`         // Date of the last release.
	Releases          int           `json:"releases"`             // Total number of releases.
	ReleasesPerDay    float64       `json:"releases_per_day"`     // Number of releases per day.
	ReleasesPerWeek   float64       `json:"releases_per_week"`    // Number of releases per week.
	ReleasesPerMonth  float64       `json:"releases_per_month"`   // Number of releases per month.
	ReleasesPerYear   float64       `json:"releases_per_year"`    // Number of releases per year.
	OpenIssuesNow     int           `json:"open_issues_now"`      // Number of currently open issues.
	OpenedIssues6M    int           `json:"opened_issues_6m"`     // Number of issues opened in the last 6 months.
	ClosedIssues6M    int           `json:"closed_issues_6m"`     // Number of issues closed in the last 6 months.
	OpenedPRs6M       int           `json:"opened_prs_6m"`        // Number of pull requests opened in the last 6 months.
	ClosedPRs6M       int           `json:"closed_prs_6m"`        // Number of pull requests closed in the last 6 months.
	Likes             int           `json:"likes"`                // Number of likes/hearts/stars.
	Forks             int           `json:"forks"`                // Number of forks.
	TopCommitters     []Contributor `json:"top_committers"`       // Top committers overall.
	TopCommitters1Y   []Contributor `json:"top_committers_1y"`    // Top committers in the last year.
}

// Contributor represents a contributor to a repository
type Contributor struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Commits int    `json:"commits"`
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
	return strings.Contains(login, "[bot]")
}

type issuePRCounts struct {
	openedIssues6M int
	closedIssues6M int
	openedPulls6M  int
	closedPulls6M  int
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
