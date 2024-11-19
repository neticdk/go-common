package repometrics

import (
	"math"
	"sort"
	"strings"
	"time"
)

type Stats struct {
	LastCommit        *time.Time    `json:"last_commit"`
	CommitsPerMonth6M int           `json:"commits_per_month_6m"`
	Contributors1Y    int           `json:"contributors_1y"`
	FirstRelease      *time.Time    `json:"first_release"`
	LastRelease       *time.Time    `json:"last_release"`
	Releases          int           `json:"releases"`
	ReleasesPerDay    float64       `json:"releases_per_day"`
	ReleasesPerWeek   float64       `json:"releases_per_week"`
	ReleasesPerMonth  float64       `json:"releases_per_month"`
	ReleasesPerYear   float64       `json:"releases_per_year"`
	OpenIssuesNow     int           `json:"open_issues_now"`
	OpenedIssues6M    int           `json:"opened_issues_6m"`
	ClosedIssues6M    int           `json:"closed_issues_6m"`
	OpenedPRs6M       int           `json:"opened_prs_6m"`
	ClosedPRs6M       int           `json:"closed_prs_6m"`
	Likes             int           `json:"likes"`
	Forks             int           `json:"forks"`
	TopCommitters     []Contributor `json:"top_committers"`
	TopCommitters1Y   []Contributor `json:"top_committers_1y"`
}

type Contributor struct {
	Name    string `json:"name"`
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
