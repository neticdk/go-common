package repometrics

import (
	"fmt"

	"github.com/neticdk/go-common/pkg/types"
)

type Scanner interface {
	Scan() ([]types.Vulnerability, error)
}

func (m *Metrics) ScanVulnerabilities(s Scanner) error {
	var err error

	vulns, err := s.Scan()
	if err != nil {
		return fmt.Errorf("failed to scan vulnerabilities: %w", err)
	}
	m.Vulnerabilities = vulns

	return nil
}
