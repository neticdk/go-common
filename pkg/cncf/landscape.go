package cncf

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	LandscapeURL = "https://raw.githubusercontent.com/cncf/landscape/master/landscape.yml"
	CacheTTL     = 1 * time.Hour
)

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
	Name        string `yaml:"name"`
	HomepageURL string `yaml:"homepage_url"`
	RepoURL     string `yaml:"repo_url"`
	Project     string `yaml:"project"`
}

// HTTPClient is an interface that abstracts the http.Get function.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// DefaultHTTPClient is the default implementation of HTTPClient that uses http.Get.
type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Get(rawURL string) (*http.Response, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	return http.Get(parsedURL.String())
}

// GetLandscape fetches the CNCF landscape from the official repository
func GetLandscape(client HTTPClient) (*Landscape, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Check if the cache is still valid
	if cache.data != nil && time.Now().Before(cache.expiresAt) {
		return cache.data, nil
	}

	res, err := client.Get(LandscapeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch landscape: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		l := &Landscape{}
		if err := yaml.Unmarshal(b, l); err != nil {
			return nil, fmt.Errorf("failed to unmarshal landscape: %w", err)
		}

		// Update the cache
		cache.data = l
		cache.expiresAt = time.Now().Add(CacheTTL)

		return l, nil
	}
	return nil, nil
}

type FindProjectOptions struct {
	RepoURL     string
	HomepageURL string
	Name        string
}

// FindProject finds a project in the landscape by repo URL, homepage URL, or name.
func (l *Landscape) FindProject(opts FindProjectOptions) *Project {
	for _, category := range l.Categories {
		for _, subcategory := range category.Subcategories {
			for _, item := range subcategory.Items {
				if item.RepoURL == opts.RepoURL || item.HomepageURL == opts.HomepageURL {
					return &item
				}
				if strings.EqualFold(item.Name, opts.Name) {
					return &item
				}
			}
		}
	}
	return nil
}
