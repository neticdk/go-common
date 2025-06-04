package cncf

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const LandscapeURL = "https://raw.githubusercontent.com/cncf/landscape/master/landscape.yml"

var CacheTTL = 1 * time.Hour

type Cache struct {
	data      *Landscape
	expiresAt time.Time
	mu        sync.Mutex
}

var cache = &Cache{}

type Landscape struct {
	Categories []Category `yaml:"landscape"`
}

type Category struct {
	Name          string        `yaml:"name"`
	Subcategories []Subcategory `yaml:"subcategories"`
}

type Subcategory struct {
	Name  string    `yaml:"name"`
	Items []Project `yaml:"items"`
}

type Project struct {
	// Native fields.
	Name        string `yaml:"name"`
	HomepageURL string `yaml:"homepage_url"`
	RepoURL     string `yaml:"repo_url"`
	Project     string `yaml:"project"`
	Crunchbase  string `yaml:"crunchbase"`
	Extra       Extra  `yaml:"extra,omitempty"`

	// The fields are added for convenience to avoid having to look them up in
	// the subcategory.
	Category    string `yaml:"category"`
	Subcategory string `yaml:"subcategory"`
}

// Extra holds extra information about the project status.
type Extra struct {
	Accepted   string  `yaml:"accepted,omitempty"`
	Incubating string  `yaml:"incubating,omitempty"`
	Graduated  string  `yaml:"graduated,omitempty"`
	Archived   string  `yaml:"archived,omitempty"`
	Audits     []Audit `yaml:"audits,omitempty"`
}

// Audit holds information about the audit of a project.
type Audit struct {
	Date   string `yaml:"date,omitempty"`
	Type   string `yaml:"type,omitempty"`
	URL    string `yaml:"url,omitempty"`
	Vendor string `yaml:"vendor,omitempty"`
}

// HTTPClient is an interface that abstracts the http.Get function.
type HTTPClient interface {
	// Get gets the content of the URL and returns the response.
	Get(ctx context.Context, rawURL string) (*http.Response, error)
}

// DefaultHTTPClient is the default implementation of HTTPClient that uses http.Get.
type DefaultHTTPClient struct{}

// Get gets the content of the URL and returns the response.
func (c *DefaultHTTPClient) Get(ctx context.Context, rawURL string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	return client.Do(req)
}

// GetLandscape fetches the CNCF landscape from the official repository
func GetLandscape(ctx context.Context, client HTTPClient) (*Landscape, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Check if the cache is still valid
	if cache.data != nil && time.Now().Before(cache.expiresAt) {
		return cache.data, nil
	}

	res, err := client.Get(ctx, LandscapeURL)
	if err != nil {
		return nil, fmt.Errorf("fetching landscape: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		l := &Landscape{}
		if err := yaml.Unmarshal(b, l); err != nil {
			return nil, fmt.Errorf("unmarshaling landscape: %w", err)
		}

		// Update the cache
		cache.data = l
		cache.expiresAt = time.Now().Add(CacheTTL)

		return l, nil
	}
	return nil, nil
}

// FindProjectOptions is a set of options to find a project in the landscape.
type FindProjectOptions struct {
	// RepoURL is the URL of the repository.
	RepoURL string

	// HomepageURL is the URL of the project's homepage.
	HomepageURL string

	// Name is the name of the project.
	Name string
}

// FindProject finds a project in the landscape by repo URL, homepage URL, or name.
func (l *Landscape) FindProject(opts FindProjectOptions) *Project {
	for _, category := range l.Categories {
		for _, subcategory := range category.Subcategories {
			for _, item := range subcategory.Items {
				foundByRepo := opts.RepoURL != "" && item.RepoURL == opts.RepoURL
				foundByHomepage := opts.HomepageURL != "" && item.HomepageURL == opts.HomepageURL
				foundByName := opts.Name != "" && strings.EqualFold(item.Name, opts.Name)
				if foundByRepo || foundByHomepage || foundByName {
					fillProject(&item, category.Name, subcategory.Name)
					return &item
				}
			}
		}
	}
	return nil
}

func fillProject(p *Project, category, subcategory string) {
	p.Category = category
	p.Subcategory = subcategory
}

// ClearCache clears the cache of the landscape data.
func ClearCache() {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.data = nil
	cache.expiresAt = time.Time{}
}
