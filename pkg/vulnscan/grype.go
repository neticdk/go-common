package vulnscan

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"sync"

	"github.com/anchore/go-version"
	"github.com/anchore/grype/grype"
	grypedb "github.com/anchore/grype/grype/db/legacy/distribution"
	grypeMatch "github.com/anchore/grype/grype/match"
	grypePkg "github.com/anchore/grype/grype/pkg"
	grypeStore "github.com/anchore/grype/grype/store"
	grypeVulnerability "github.com/anchore/grype/grype/vulnerability"
	syftSbom "github.com/anchore/syft/syft/sbom"
	"github.com/neticdk/go-common/pkg/sbom"
	"github.com/neticdk/go-common/pkg/types"
)

var grypeDBRootDir = "/tmp/grype-db"

// GrypeScanner is a scanner that uses Grype to find vulnerabilities in a
// project
type GrypeScanner struct {
	manifestScanner
}

// NewGrypeScanner creates a new GrypeScanner
func NewGrypeScanner(path string) *GrypeScanner {
	return &GrypeScanner{
		manifestScanner{
			path: path,
		},
	}
}

// Scan scans the project in the given path and returns a list of vulnerabilities
func (s *GrypeScanner) Scan(ctx context.Context) ([]types.Vulnerability, error) {
	logger := slog.Default()
	var err error
	if s.path == "" {
		s.path, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
	}

	sboms, err := sbom.GenerateSBOMsFromPath(ctx, s.path)
	if err != nil {
		return nil, fmt.Errorf("failed to generate SBOM: %w", err)
	}

	vulnerabilities := []types.Vulnerability{}
	if len(sboms) == 0 {
		return vulnerabilities, nil
	}

	for _, s := range sboms {
		vulns, err := GrypeScanSBOM(ctx, *s)
		if err != nil {
			logger.WarnContext(ctx, "failed to get vulnerabilities from SBOM", slog.Any("error", err))
			continue
		}
		vulnerabilities = append(vulnerabilities, vulns...)
	}
	return vulnerabilities, nil
}

// GrypeScanSBOM extracts vulnerabilities from the given SBOM.
// It loads the Grype vulnerability database, matches the packages in the SBOM
// against known vulnerabilities, and returns a list of vulnerabilities.
func GrypeScanSBOM(ctx context.Context, s syftSbom.SBOM) ([]types.Vulnerability, error) {
	vulns := []types.Vulnerability{}

	dbConfig := grypedb.Config{
		DBRootDir:  grypeDBRootDir,
		ListingURL: "https://toolbox-data.anchore.io/grype/databases/listing.json",
	}
	datastore, _, _, err := grype.LoadVulnerabilityDB(dbConfig, true)
	if err != nil {
		return nil, fmt.Errorf("failed to load vulnerability database: %w", err)
	}

	matcher := grype.DefaultVulnerabilityMatcher(*datastore)

	syftPkgs := s.Artifacts.Packages.Sorted()
	grypePkgs := grypePkg.FromPackages(syftPkgs, grypePkg.SynthesisConfig{GenerateMissingCPEs: false})

	col, _, err := matcher.FindMatches(grypePkgs, grypePkg.Context{
		Source: &s.Source,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find matches: %w", err)
	}

	// Create channels for work distribution and result collection
	workChan := make(chan grypeMatch.Match, len(col.Sorted()))
	resultChan := make(chan types.Vulnerability, len(col.Sorted()))

	// Use a WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Number of workers
	numWorkers := 10

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for match := range workChan {
				grypeProcessMatch(ctx, match, datastore, resultChan)
			}
		}()
	}

	// Send work to the work channel
	for _, match := range col.Sorted() {
		workChan <- match
	}
	close(workChan)

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Collect results from the result channel
	for vuln := range resultChan {
		vulns = append(vulns, vuln)
	}

	return vulns, nil
}

func grypeProcessMatch(ctx context.Context, match grypeMatch.Match, datastore *grypeStore.Store, resultChan chan<- types.Vulnerability) {
	logger := slog.Default()
	metadata, err := datastore.GetMetadata(match.Vulnerability.ID, match.Vulnerability.Namespace)
	if err != nil {
		logger.WarnContext(ctx, "getting metadata for vulnerability", slog.String("ID", match.Vulnerability.ID), slog.Any("error", err))
		return
	}

	// Only interested in the highest CVSS version
	versions := metadata.Cvss
	sort.Sort(grypeByCVSSVersion(versions))
	var cvss *types.CVSS
	if len(versions) > 0 {
		cvss = &types.CVSS{
			Vector:    versions[len(versions)-1].Vector,
			BaseScore: versions[len(versions)-1].Metrics.BaseScore,
		}
	}
	resultChan <- types.Vulnerability{
		ID:          match.Vulnerability.ID,
		Severity:    metadata.Severity,
		PackageName: match.Package.Name,
		Description: metadata.Description,
		FixState:    string(match.Vulnerability.Fix.State),
		CVSS:        cvss,
	}
}

func GrypeSetDBRootDir(dir string) {
	grypeDBRootDir = dir
}

type grypeByCVSSVersion []grypeVulnerability.Cvss

func (s grypeByCVSSVersion) Len() int {
	return len(s)
}

func (s grypeByCVSSVersion) Swap(i, j int) {
	s[i].Version, s[j].Version = s[j].Version, s[i].Version
}

func (s grypeByCVSSVersion) Less(i, j int) bool {
	logger := slog.Default()
	v1, err1 := version.NewVersion(s[i].Version)
	v2, err2 := version.NewVersion(s[j].Version)

	if err1 != nil || err2 != nil {
		// Treat versions as equal if parsing fails
		logger.DebugContext(context.TODO(), "parsing version", slog.Any("error1", err1), slog.Any("error2", err2))
		return false
	}

	return v1.LessThan(v2)
}
