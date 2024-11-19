# CNCF Package

This package provides functionality to interact with the [CNCF Landscape](https://landscape.cncf.io/),
specifically for querying project information from the landscape's data.

## Features

- Fetch and parse the official CNCF landscape data
- Cache landscape data with configurable TTL
- Search for projects by:
  - Repository URL
  - Homepage URL
  - Project name (case-insensitive)

## Usage

```go
import "github.com/neticdk/go-common/pkg/cncf"

// Create a client (or use your own that implements HTTPClient)
client := &cncf.DefaultHTTPClient{}

// Get the landscape data (cached for 1 hour by default)
landscape, err := cncf.GetLandscape(client)
if err != nil {
    // Handle error
}

// Find a project by repo URL
project := landscape.FindProject(cncf.FindProjectOptions{
    RepoURL: "https://github.com/kubernetes/kubernetes",
})

// Find a project by name
project = landscape.FindProject(cncf.FindProjectOptions{
    Name: "Kubernetes",
})
```

## Cache Behavior

The package includes built-in caching of landscape data to avoid unnecessary network requests.
The default cache TTL is 1 hour but can be modified by changing the `CacheTTL` constant.
