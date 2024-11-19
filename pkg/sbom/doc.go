/*
Package sbom provides functionality for working with Software Bill of Materials (SBOM) files.

This package supports:

- generating SBOMs from parsing kubernetes manifest files.

Basic usage:

	// Parse an SBOM file
	sbom, err := sbom.Parse("path/to/sbom.json")
	if err != nil {
	    // Handle error
	}

	// Get components
	components := sbom.Components()

	// Check licenses
	licenses := sbom.Licenses()

The package supports different SBOM formats and provides utilities for:
  - Parsing SBOM files
  - Extracting component information
  - License analysis
  - Dependency tracking
*/
package sbom
