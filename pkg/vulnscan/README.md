# Vulnerability Scanning Package

This package provides comprehensive vulnerability scanning capabilities for Go
applications. It supports many scanning backends and vulnerability databases.

## Usage

```go
import "github.com/neticdk/go-common/pkg/vulnscan"

// Create a new scanner
scanner := vulnscan.NewScanner(vulnscan.Options{
    // Configure scanner options
})

// Perform a scan
results, err := scanner.Scan(context.Background())
if err != nil {
    // Handle error
}

// Process results
for _, vuln := range results.Vulnerabilities {
    fmt.Printf("Found vulnerability: %s (severity: %s)\n", vuln.ID, vuln.Severity)
}
```

## Configuration

A scanner implements the `Scanner` interface. Specifically, the `Scan` method
which returns a list of vulnerabilities and an error.

For now, GrypeScanner is the only implementation.
