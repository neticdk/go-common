package repometrics

import (
	"context"
	"fmt"

	"github.com/neticdk/go-common/pkg/types"
)

// Scanner is the interface for vulnerability scanners
type Scanner interface {
	// Scan scans for vulnerabilities
	Scan(context.Context) ([]types.Vulnerability, error)
}

// ScanVulnerabilities scans for vulnerabilities and updates the metrics
func (m *Metrics) ScanVulnerabilities(ctx context.Context, s Scanner) error {
	var err error

	vulnerabilities, err := s.Scan(ctx)
	if err != nil {
		return fmt.Errorf("scanning vulnerabilities: %w", err)
	}
	// m.Vulnerabilities = vulnerabilities
	// produce a cycloneDX file from scanning, whether it is done by SBOM or image or File
	fmt.Println("Vulnerabilities should be produced in CycloneDX format")
	fmt.Println(vulnerabilities)
	return nil
}
