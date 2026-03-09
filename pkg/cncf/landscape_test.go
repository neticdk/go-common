package cncf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	Response *http.Response
	Err      error
	Context  context.Context
}

func (m *MockHTTPClient) Get(ctx context.Context, rawUrl string) (*http.Response, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
	}

	return m.Response, m.Err
}

func TestGetLandscape(t *testing.T) {
	mockData := `
landscape:
  - category:
    name: Provisioning
    subcategories:
      - subcategory:
        name: Automation & Configuration
        items:
          - item:
            name: Airship
            homepage_url: https://www.airshipit.org/
            repo_url: https://github.com/airshipit/treasuremap
            project: airship
            crunchbase: https://www.crunchbase.com/organization/cloud-native-computing-foundation
            extra:
              accepted: '2022-09-14'
              incubating: '2022-09-14'
              archived: '2022-09-14'
              graduated: '2024-09-11'
              audits:
                - date: '2024-01-01'
                  type: security
                  url: https://github.com/airshipit/treasuremap/blob/main/audits/security.md
                  vendor: Cure53
`

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(mockData)),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
		Context:  t.Context(),
	}

	// Clear the cache before running the test
	cache = &Cache{}

	landscape, err := GetLandscape(t.Context(), mockClient)
	assert.NoError(t, err)
	assert.NotNil(t, landscape)
	assert.Equal(t, "Provisioning", landscape.Categories[0].Name)
	assert.Equal(t, "Automation & Configuration", landscape.Categories[0].Subcategories[0].Name)
	assert.Equal(t, "Airship", landscape.Categories[0].Subcategories[0].Items[0].Name)
	assert.Equal(t, "https://www.airshipit.org/", landscape.Categories[0].Subcategories[0].Items[0].HomepageURL)
	assert.Equal(t, "https://github.com/airshipit/treasuremap", landscape.Categories[0].Subcategories[0].Items[0].RepoURL)
	assert.Equal(t, "https://www.crunchbase.com/organization/cloud-native-computing-foundation", landscape.Categories[0].Subcategories[0].Items[0].Crunchbase)
	assert.Equal(t, "airship", landscape.Categories[0].Subcategories[0].Items[0].Project)
	assert.Equal(t, "2022-09-14", landscape.Categories[0].Subcategories[0].Items[0].Extra.Accepted)
	assert.Equal(t, "2022-09-14", landscape.Categories[0].Subcategories[0].Items[0].Extra.Incubating)
	assert.Equal(t, "2024-09-11", landscape.Categories[0].Subcategories[0].Items[0].Extra.Graduated)
	assert.Equal(t, "2022-09-14", landscape.Categories[0].Subcategories[0].Items[0].Extra.Archived)
	assert.Equal(t, 1, len(landscape.Categories[0].Subcategories[0].Items[0].Extra.Audits))
	assert.Equal(t, "security", landscape.Categories[0].Subcategories[0].Items[0].Extra.Audits[0].Type)
}

func TestGetLandscape_Cache(t *testing.T) {
	mockData := `
landscape:
  - category:
    name: Provisioning
    subcategories:
      - subcategory:
        name: Automation & Configuration
        items:
          - item:
            name: Airship
            homepage_url: https://www.airshipit.org/
            repo_url: https://github.com/airshipit/treasuremap
            project: airship
            crunchbase: https://www.crunchbase.com/organization/cloud-native-computing-foundation
            extra:
              accepted: '2022-09-14'
              incubating: '2022-09-14'
              archived: '2022-09-14'
              graduated: '2024-09-11'
              audits:
                - date: '2024-01-01'
                  type: security
                  url: https://github.com/airshipit/treasuremap/blob/main/audits/security.md
                  vendor: Cure53
`

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(mockData)),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}

	// Clear the cache before running the test
	cache = &Cache{}

	// First call to populate the cache
	_, err := GetLandscape(t.Context(), mockClient)
	assert.NoError(t, err)

	// Modify the mock client to return an error
	mockClient.Err = fmt.Errorf("network error")

	// Second call should return cached data
	landscape, err := GetLandscape(t.Context(), mockClient)
	assert.NoError(t, err)
	assert.NotNil(t, landscape)
	assert.Equal(t, "Provisioning", landscape.Categories[0].Name)
	assert.Equal(t, "Automation & Configuration", landscape.Categories[0].Subcategories[0].Name)
	assert.Equal(t, "Airship", landscape.Categories[0].Subcategories[0].Items[0].Name)
	assert.Equal(t, "https://www.airshipit.org/", landscape.Categories[0].Subcategories[0].Items[0].HomepageURL)
	assert.Equal(t, "https://github.com/airshipit/treasuremap", landscape.Categories[0].Subcategories[0].Items[0].RepoURL)
	assert.Equal(t, "airship", landscape.Categories[0].Subcategories[0].Items[0].Project)
	assert.Equal(t, "2022-09-14", landscape.Categories[0].Subcategories[0].Items[0].Extra.Accepted)
	assert.Equal(t, "2022-09-14", landscape.Categories[0].Subcategories[0].Items[0].Extra.Incubating)
	assert.Equal(t, "2024-09-11", landscape.Categories[0].Subcategories[0].Items[0].Extra.Graduated)
	assert.Equal(t, "2022-09-14", landscape.Categories[0].Subcategories[0].Items[0].Extra.Archived)
	assert.Equal(t, 1, len(landscape.Categories[0].Subcategories[0].Items[0].Extra.Audits))
	assert.Equal(t, "security", landscape.Categories[0].Subcategories[0].Items[0].Extra.Audits[0].Type)
	assert.Equal(t, 1, len(cache.data.Categories)) // Ensure item is cached
}

func TestLandscape_FindProject(t *testing.T) {
	landscape := &Landscape{
		Categories: []Category{
			{
				Name: "Category 1",
				Subcategories: []Subcategory{
					{
						Name: "Subcategory 1",
						Items: []Project{
							{
								Name:        "Project 1",
								HomepageURL: "https://project1.com",
								RepoURL:     "https://github.com/org/project1",
							},
							{
								Name:        "Project 2",
								HomepageURL: "https://project2.com",
								RepoURL:     "https://github.com/org/project2",
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name    string
		opts    FindProjectOptions
		want    *Project
		wantNil bool
	}{
		{
			name: "find by repo URL",
			opts: FindProjectOptions{
				RepoURL: "https://github.com/org/project1",
			},
			want: &Project{
				Name:        "Project 1",
				HomepageURL: "https://project1.com",
				RepoURL:     "https://github.com/org/project1",
			},
		},
		{
			name: "find by homepage URL",
			opts: FindProjectOptions{
				HomepageURL: "https://project2.com",
			},
			want: &Project{
				Name:        "Project 2",
				HomepageURL: "https://project2.com",
				RepoURL:     "https://github.com/org/project2",
			},
		},
		{
			name: "find by name - exact match",
			opts: FindProjectOptions{
				Name: "Project 1",
			},
			want: &Project{
				Name:        "Project 1",
				HomepageURL: "https://project1.com",
				RepoURL:     "https://github.com/org/project1",
			},
		},
		{
			name: "find by name - case insensitive",
			opts: FindProjectOptions{
				Name: "project 2",
			},
			want: &Project{
				Name:        "Project 2",
				HomepageURL: "https://project2.com",
				RepoURL:     "https://github.com/org/project2",
			},
		},
		{
			name: "not found - non-existent repo URL",
			opts: FindProjectOptions{
				RepoURL: "https://github.com/org/nonexistent",
			},
			wantNil: true,
		},
		{
			name: "not found - non-existent homepage URL",
			opts: FindProjectOptions{
				HomepageURL: "https://nonexistent.com",
			},
			wantNil: true,
		},
		{
			name: "not found - non-existent name",
			opts: FindProjectOptions{
				Name: "Nonexistent Project",
			},
			wantNil: true,
		},
		{
			name:    "not found - empty options",
			opts:    FindProjectOptions{},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := landscape.FindProject(tt.opts)

			if tt.wantNil {
				if got != nil {
					t.Errorf("FindProject() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("FindProject() returned nil, want non-nil")
			}

			if got.Name != tt.want.Name {
				t.Errorf("FindProject().Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.HomepageURL != tt.want.HomepageURL {
				t.Errorf("FindProject().HomepageURL = %v, want %v", got.HomepageURL, tt.want.HomepageURL)
			}
			if got.RepoURL != tt.want.RepoURL {
				t.Errorf("FindProject().RepoURL = %v, want %v", got.RepoURL, tt.want.RepoURL)
			}
		})
	}
}

func TestClearCache(t *testing.T) {
	// Set some dummy values to ensure they are cleared
	cache.data = &Landscape{} // Use a dummy non-nil value
	cache.expiresAt = time.Now().Add(1 * time.Hour)

	// Call the function to clear the cache
	ClearCache()

	// Assert that the cache is cleared
	if cache.data != nil {
		t.Errorf("cache.data was not cleared, expected nil, got %+v", cache.data)
	}

	var zeroTime time.Time
	if !cache.expiresAt.Equal(zeroTime) {
		t.Errorf("cache.expiresAt was not cleared, expected zero time (%s), got %s", zeroTime, cache.expiresAt)
	}
}
