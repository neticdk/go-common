/*
Package vulnscan provides vulnerability scanning capabilities for software projects.

This package implements the Scanner interface to allow using various
vulnerability scanning tools and databases, allowing for comprehensive security
analysis of software components. It provides a unified interface for
vulnerability detection and reporting.

Basic usage:

		scanner := vulnscan.NewScanner(vulnscan.Options{
	        // Configure scanner options
	    })

		// Scan a project
		results, err := scanner.Scan(context.Background())
		if err != nil {
		    // Handle error
		}

		// Process results
		for _, vuln := range results.Vulnerabilities {
		    // Handle each vulnerability
		}

The package supports:
  - Multiple vulnerability databases
  - Various scanning backends
  - Detailed vulnerability reporting
*/
package vulnscan
