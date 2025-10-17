package repometrics

import (
	"testing"
	"time"
)

func TestCalculateReleaseMetrics(t *testing.T) {
	now := time.Now()
	oneYearAgo := now.Add(-365.25 * 24 * time.Hour)
	sixMonthsAgo := now.Add(-182.625 * 24 * time.Hour)

	tests := []struct {
		name         string
		firstRelease *time.Time
		lastRelease  *time.Time
		releases     int
		want         ReleaseMetrics
	}{
		{
			name:         "nil dates",
			firstRelease: nil,
			lastRelease:  nil,
			releases:     10,
			want:         ReleaseMetrics{},
		},
		{
			name:         "zero releases",
			firstRelease: &oneYearAgo,
			lastRelease:  &now,
			releases:     0,
			want:         ReleaseMetrics{},
		},
		{
			name:         "negative time delta",
			firstRelease: &now,
			lastRelease:  &oneYearAgo,
			releases:     0,
			want:         ReleaseMetrics{},
		},
		{
			name:         "one year period",
			firstRelease: &oneYearAgo,
			lastRelease:  &now,
			releases:     12,
			want: ReleaseMetrics{
				PerDay:   0.0,
				PerWeek:  0.2,
				PerMonth: 1.0,
				PerYear:  12.0,
			},
		},
		{
			name:         "six month period",
			firstRelease: &sixMonthsAgo,
			lastRelease:  &now,
			releases:     6,
			want: ReleaseMetrics{
				PerDay:   0.0,
				PerWeek:  0.2,
				PerMonth: 1.0,
				PerYear:  12.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateReleaseMetrics(tt.firstRelease, tt.lastRelease, tt.releases)

			if got.PerDay != tt.want.PerDay {
				t.Errorf("PerDay = %v, want %v", got.PerDay, tt.want.PerDay)
			}
			if got.PerWeek != tt.want.PerWeek {
				t.Errorf("PerWeek = %v, want %v", got.PerWeek, tt.want.PerWeek)
			}
			if got.PerMonth != tt.want.PerMonth {
				t.Errorf("PerMonth = %v, want %v", got.PerMonth, tt.want.PerMonth)
			}
			if got.PerYear != tt.want.PerYear {
				t.Errorf("PerYear = %v, want %v", got.PerYear, tt.want.PerYear)
			}
		})
	}
}

func TestSortContributors(t *testing.T) {
	tests := []struct {
		name         string
		contributors []Contributor
		want         []Contributor
	}{
		{
			name:         "empty list",
			contributors: []Contributor{},
			want:         []Contributor{},
		},
		{
			name: "single contributor",
			contributors: []Contributor{
				{Name: "Abbey", Commits: 10},
			},
			want: []Contributor{
				{Name: "Abbey", Commits: 10},
			},
		},
		{
			name: "already sorted",
			contributors: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
			},
			want: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
			},
		},
		{
			name: "unsorted list",
			contributors: []Contributor{
				{Name: "Raffo", Commits: 5},
				{Name: "Abbey", Commits: 10},
			},
			want: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
			},
		},
		{
			name: "contributors with equal commits",
			contributors: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 10},
			},
			want: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortContributors(tt.contributors)

			if len(tt.contributors) != len(tt.want) {
				t.Fatalf("got %v contributors, want %v", len(tt.contributors), len(tt.want))
			}

			for i := range tt.contributors {
				if tt.contributors[i] != tt.want[i] {
					t.Errorf("at index %d, got %v, want %v", i, tt.contributors[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetTopContributors(t *testing.T) {
	tests := []struct {
		name         string
		contributors []Contributor
		limit        int
		want         []Contributor
	}{
		{
			name:         "empty list",
			contributors: []Contributor{},
			limit:        5,
			want:         []Contributor{},
		},
		{
			name: "limit greater than contributors",
			contributors: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
			},
			limit: 5,
			want: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
			},
		},
		{
			name: "limit less than contributors",
			contributors: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
				{Name: "Tziki", Commits: 3},
			},
			limit: 2,
			want: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
			},
		},
		{
			name: "limit equal to contributors",
			contributors: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
			},
			limit: 2,
			want: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
			},
		},
		{
			name: "zero limit",
			contributors: []Contributor{
				{Name: "Abbey", Commits: 10},
				{Name: "Raffo", Commits: 5},
			},
			limit: 0,
			want:  []Contributor{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTopContributors(tt.contributors, tt.limit)

			if len(got) != len(tt.want) {
				t.Fatalf("got %v contributors, want %v", len(got), len(tt.want))
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("at index %d, got %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestIsBot(t *testing.T) {
	tests := []struct {
		name  string
		login string
		want  bool
	}{
		{
			name:  "login contains [bot]",
			login: "example[bot]",
			want:  true,
		},
		{
			name:  "login contains bot",
			login: "examplebot",
			want:  true,
		},
		{
			name:  "login contains istio-testing",
			login: "istio-testing",
			want:  true,
		},
		{
			name:  "login contains fluxcdbot",
			login: "fluxcdbot",
			want:  true,
		},
		{
			name:  "login contains dependabot",
			login: "dependabot",
			want:  true,
		},
		{
			name:  "login contains renovate-bot",
			login: "renovate-bot",
			want:  true,
		},
		{
			name:  "login contains renovate[bot]",
			login: "renovate[bot]",
			want:  true,
		},
		{
			name:  "login contains action",
			login: "github-action",
			want:  true,
		},
		{
			name:  "login does not contain bot-related keywords",
			login: "regular-user",
			want:  false,
		},
		{
			name:  "empty login",
			login: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBot(tt.login)
			if got != tt.want {
				t.Errorf("isBot(%q) = %v, want %v", tt.login, got, tt.want)
			}
		})
	}
}
