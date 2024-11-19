# SBOM Package

This package provides functionality for working with Software Bill of Materials (SBOM) files in Go applications.

## Features

- Parse SBOM files in multiple formats:
  - CycloneDX
  - SPDX
- Extract component information
- Analyze licenses
- Track dependencies
- Validate SBOM structure

## Usage

```go
import "github.com/neticdk/go-common/pkg/sbom"

// Parse an SBOM file
sbom, err := sbom.Parse("path/to/sbom.json")
if err != nil {
    // Handle error
}

// Get all components
components := sbom.Components()

// Get all licenses
licenses := sbom.Licenses()

// Validate SBOM
if err := sbom.Validate(); err != nil {
    // Handle validation errors
}
```

## Supported Formats

The package currently supports the following SBOM formats:
- CycloneDX JSON
- CycloneDX XML
- SPDX JSON
- SPDX Tag-Value

## Integration

This package is designed to work seamlessly with various SBOM generation tools and can be integrated into existing workflows for software composition analysis and compliance checking.