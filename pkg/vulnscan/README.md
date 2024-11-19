# Vulnerability Scanning Package

This package provides comprehensive vulnerability scanning capabilities for Go
applications, supporting multiple scanning backends and vulnerability databases.

## Usage

```go
import "github.com/neticdk/go-common/pkg/vulnscan"

// Create a new scanner
scanner := vulnscan.NewScanner(vulnscan.Options{
    // Configure scanner options
})

// Perform a scan
results, err := scanner.Scan()
if err != nil {
    // Handle error
}

// Process results
for _, vuln := range results.Vulnerabilities {
    fmt.Printf("Found vulnerability: %s (severity: %s)\n", vuln.ID, vuln.Severity)
}
```

## Configuration

The scanner must implement the `Scanner` interface, which includes a `Scan`
method that returns a list of vulnerabilities and an error.

For now the only implementation is the GrypeScanner, which uses the Grype.
