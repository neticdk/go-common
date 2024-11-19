package repometrics

import (
	"fmt"

	"github.com/neticdk/go-common/pkg/types"
)

// Scanner is the interface for vulnerability scanners
type Scanner interface {
	// Scan scans for vulnerabilities
	Scan() ([]types.Vulnerability, error)
}

// ScanVulnerabilities scans for vulnerabilities and updates the metrics
func (m *Metrics) ScanVulnerabilities(s Scanner) error {
	var err error

	vulns, err := s.Scan()
	if err != nil {
		return fmt.Errorf("failed to scan vulnerabilities: %w", err)
	}
	m.Vulnerabilities = vulns

	return nil
}
