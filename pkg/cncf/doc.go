/*
Package cncf provides functionality for interacting with the Cloud Native Computing Foundation (CNCF) landscape.

The package allows fetching and querying the official CNCF landscape data, which contains information about
all CNCF projects, their status, and related metadata. It includes built-in caching to minimize network requests
and provides a simple interface for finding projects by various identifiers.

Basic usage:

	client := &cncf.DefaultHTTPClient{}
	landscape, err := cncf.GetLandscape(client)
	if err != nil {
		// Handle error
	}

	// Find a project
	project := landscape.FindProject(cncf.FindProjectOptions{
		Name: "Kubernetes",
	})

The landscape data is automatically cached for one hour (configurable via CacheTTL) to reduce
network requests and improve performance.

Projects can be found by:
  - Repository URL
  - Homepage URL
  - Project name (case-insensitive matching)
*/
package cncf
