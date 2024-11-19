package repometrics

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		typ     RepoType
		wantErr bool
	}{
		{
			name:    "github repository",
			typ:     RepoTypeGitHub,
			wantErr: false,
		},
		{
			name:    "empty repository type",
			typ:     "",
			wantErr: true,
		},
		{
			name:    "invalid repository type",
			typ:     "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.typ)

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if got != nil {
					t.Error("New() returned non-nil Metrics when error expected")
				}
				return
			}

			if got == nil {
				t.Fatal("New() returned nil when no error expected")
			}
			if got.Type != tt.typ {
				t.Errorf("New().Type = %v, want %v", got.Type, tt.typ)
			}
			if got.Stats == nil {
				t.Error("New().Stats is nil, want initialized Stats")
			}
			if got.Vulnerabilities == nil {
				t.Error("New().Vulnerabilities is nil, want initialized slice")
			}
			if len(got.Vulnerabilities) != 0 {
				t.Errorf("New().Vulnerabilities has len %d, want 0", len(got.Vulnerabilities))
			}
		})
	}
}

func TestMetrics_ProjectAge(t *testing.T) {
	tests := []struct {
		name      string
		createdAt *time.Time
		want      time.Duration
	}{
		{
			name:      "nil creation date",
			createdAt: nil,
			want:      time.Duration(0),
		},
		{
			name: "specific creation date",
			createdAt: func() *time.Time {
				t := time.Now().Add(-24 * time.Hour)
				return &t
			}(),
			want: 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				CreatedAt: tt.createdAt,
			}
			got := m.ProjectAge()

			// For nil creation date, expect exactly 0
			if tt.createdAt == nil {
				if got != 0 {
					t.Errorf("ProjectAge() = %v, want 0", got)
				}
				return
			}

			// For actual dates, allow for small timing differences
			diff := got - tt.want
			if diff < 0 {
				diff = -diff
			}
			if diff > time.Second {
				t.Errorf("ProjectAge() = %v, want ~%v (diff too large: %v)", got, tt.want, diff)
			}
		})
	}
}
