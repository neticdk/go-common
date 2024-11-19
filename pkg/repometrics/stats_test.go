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
