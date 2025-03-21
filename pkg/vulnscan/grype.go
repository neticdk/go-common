package vulnscan

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/anchore/go-version"
	"github.com/anchore/grype/grype"
	"github.com/anchore/grype/grype/db/v6/distribution"
	"github.com/anchore/grype/grype/db/v6/installation"
	grypeMatch "github.com/anchore/grype/grype/match"
	"github.com/anchore/grype/grype/matcher"
	"github.com/anchore/grype/grype/matcher/dotnet"
	"github.com/anchore/grype/grype/matcher/golang"
	"github.com/anchore/grype/grype/matcher/java"
	"github.com/anchore/grype/grype/matcher/javascript"
	"github.com/anchore/grype/grype/matcher/python"
	"github.com/anchore/grype/grype/matcher/ruby"
	"github.com/anchore/grype/grype/matcher/rust"
	"github.com/anchore/grype/grype/matcher/stock"
	grypePkg "github.com/anchore/grype/grype/pkg"
	grypeVulnerability "github.com/anchore/grype/grype/vulnerability"
	syftSbom "github.com/anchore/syft/syft/sbom"
	"github.com/neticdk/go-common/pkg/sbom"
	"github.com/neticdk/go-common/pkg/types"
)

const (
	DefaultDBRootDir = "/tmp/grypedb"

	// Number of scanner workers
	scannerWorkers = 10

	mavenSearchBaseURL = "https://search.maven.org/solrsearch/select"
)

// GrypeScanner is a scanner that uses Grype to find vulnerabilities in a
// project
type GrypeScanner struct {
	manifestScanner
	dbRootDir          string
	cleanupDBAfterScan bool
}

// GrypeScannerOptions specifies the options for the GrypeScanner
type GrypeScannerOptions struct {
	ManifestPath       string    // ManifestPath specifies the path to the project manifest
	Manifest           io.Reader // Manifest is a readable representation of a manifest file
	DBRootDir          string    // DBRootDir specifies the root directory of the Grype database
	CleanupDBAfterScan bool      // CleanupDBAfterScan specifies whether to clean up the Grype database after the scan
}

// DefaultGrypeScannerOptions returns the default GrypeScannerOptions
// It sets:
// - DBRootDir to the default Grype database root directory
// - CleanupDBAfterScan to false
func DefaultGrypeScannerOptions() *GrypeScannerOptions {
	var wd string
	wd, _ = os.Getwd()
	return &GrypeScannerOptions{
		ManifestPath:       wd,
		Manifest:           strings.NewReader(""),
		DBRootDir:          DefaultDBRootDir,
		CleanupDBAfterScan: false,
	}
}

// NewGrypeScanner creates a new GrypeScanner
func NewGrypeScanner(opts *GrypeScannerOptions) *GrypeScanner {
	if opts == nil {
		opts = DefaultGrypeScannerOptions()
	}
	s := &GrypeScanner{
		manifestScanner{
			manifestPath: opts.ManifestPath,
			manifest:     opts.Manifest,
		},
		opts.DBRootDir,
		opts.CleanupDBAfterScan,
	}
	if s.dbRootDir == "" {
		s.dbRootDir = DefaultDBRootDir
	}
	if s.manifest == nil {
		s.manifest = strings.NewReader("")
	}
	return s
}

// Scan scans the project in the given path and returns a list of vulnerabilities
func (s *GrypeScanner) Scan(ctx context.Context) ([]types.Vulnerability, error) {
	logger := slog.Default()

	sboms, err := sbom.GenerateSBOMsFromPath(ctx, s.manifestPath)
	if err != nil {
		return nil, fmt.Errorf("generating SBOM: %w", err)
	}
	sboms2, err := sbom.GenerateSBOMsFromManifest(ctx, s.manifest)
	if err != nil {
		return nil, fmt.Errorf("generating SBOM: %w", err)
	}
	sboms = append(sboms, sboms2...)

	vulnerabilities := []types.Vulnerability{}
	if len(sboms) == 0 {
		return vulnerabilities, nil
	}

	for _, sbm := range sboms {
		vulns, err := s.GrypeScanSBOM(ctx, *sbm)
		if err != nil {
			logger.WarnContext(ctx, "getting vulnerabilities from SBOM", slog.Any("error", err))
			continue
		}
		vulnerabilities = append(vulnerabilities, vulns...)
	}

	sort.SliceStable(vulnerabilities, func(i, j int) bool {
		baseScoreI := 0.0
		baseScoreJ := 0.0
		if vulnerabilities[i].CVSS != nil {
			baseScoreI = vulnerabilities[i].CVSS.BaseScore
		}
		if vulnerabilities[j].CVSS != nil {
			baseScoreJ = vulnerabilities[j].CVSS.BaseScore
		}
		return baseScoreI > baseScoreJ
	})

	return vulnerabilities, nil
}

// GrypeScanSBOM extracts vulnerabilities from the given SBOM.
// It loads the Grype vulnerability database, matches the packages in the SBOM
// against known vulnerabilities, and returns a list of vulnerabilities.
func (s *GrypeScanner) GrypeScanSBOM(ctx context.Context, sbm syftSbom.SBOM) ([]types.Vulnerability, error) {
	logger := slog.Default()
	vulns := []types.Vulnerability{}

	distConfig := distribution.Config{
		LatestURL: "https://grype.anchore.io/databases/v6/latest.json",
	}
	installCfg := installation.Config{
		DBRootDir: s.dbRootDir,
	}
	client, err := distribution.NewClient(distConfig)
	if err != nil {
		return nil, fmt.Errorf("creating distribution client: %w", err)
	}
	if s.cleanupDBAfterScan {
		dbCurator, err := installation.NewCurator(installCfg, client)
		if err != nil {
			return nil, fmt.Errorf("creating database curator: %w", err)
		}
		defer func() {
			if err := dbCurator.Delete(); err != nil {
				logger.ErrorContext(ctx, "cleaning up database", slog.Any("error", err))
			}
		}()
	}

	datastore, _, err := grype.LoadVulnerabilityDB(distConfig, installCfg, true)
	if err != nil {
		return nil, fmt.Errorf("loading vulnerability database: %w", err)
	}

	vulnMatcher := grype.VulnerabilityMatcher{
		VulnerabilityProvider: datastore,
		Matchers:              createMatchers(true),
	}

	syftPkgs := sbm.Artifacts.Packages.Sorted()
	grypePkgs := grypePkg.FromPackages(syftPkgs, grypePkg.SynthesisConfig{GenerateMissingCPEs: false})

	matches, _, err := vulnMatcher.FindMatches(grypePkgs, grypePkg.Context{
		Source: &sbm.Source,
	})
	if err != nil {
		return nil, fmt.Errorf("finding matches: %w", err)
	}

	matchCount := len(matches.Sorted())
	// Create channels for work distribution and result collection
	workChan := make(chan grypeMatch.Match, matchCount)
	resultChan := make(chan types.Vulnerability, matchCount)

	// Use a WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Start workers
	for range scannerWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for match := range workChan {
				grypeProcessMatch(ctx, match, datastore, resultChan)
			}
		}()
	}

	// Send work to the work channel
	for _, match := range matches.Sorted() {
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

func grypeProcessMatch(ctx context.Context, match grypeMatch.Match, datastore grypeVulnerability.Provider, resultChan chan<- types.Vulnerability) {
	logger := slog.Default()
	metadata, err := datastore.VulnerabilityMetadata(match.Vulnerability.Reference)
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

func createMatchers(useCPEs bool) []grypeMatch.Matcher {
	return matcher.NewDefaultMatchers(
		matcher.Config{
			Dotnet: dotnet.MatcherConfig{UseCPEs: useCPEs},
			Golang: golang.MatcherConfig{
				UseCPEs:                                true,
				AlwaysUseCPEForStdlib:                  true,
				AllowMainModulePseudoVersionComparison: false,
			},
			Java: java.MatcherConfig{
				ExternalSearchConfig: java.ExternalSearchConfig{
					SearchMavenUpstream: true,
					MavenBaseURL:        mavenSearchBaseURL,
					MavenRateLimit:      400 * time.Millisecond, // increased from the default of 300ms to avoid rate limiting with extremely large set of java packages such as druid
				},
				UseCPEs: useCPEs,
			},
			Javascript: javascript.MatcherConfig{UseCPEs: useCPEs},
			Python:     python.MatcherConfig{UseCPEs: useCPEs},
			Ruby:       ruby.MatcherConfig{UseCPEs: useCPEs},
			Rust:       rust.MatcherConfig{UseCPEs: useCPEs},
			Stock:      stock.MatcherConfig{UseCPEs: true},
		},
	)
}
