package helm

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"helm.sh/helm/v4/pkg/registry"
)

func TestPullChart(t *testing.T) {
	testdataDir := filepath.Join("testdata")

	// Create a temporary directory for test data
	tmpDir, err := os.MkdirTemp("", "go-common-helm-test-")
	if err != nil {
		t.Fatalf("creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Start a local HTTP server to serve the chart tarball
	httpServer := httptest.NewServer(http.FileServer(http.Dir(testdataDir)))
	defer httpServer.Close()

	// Test pulling the chart from the HTTP server
	t.Run("HTTP server", func(t *testing.T) {
		dstDir := filepath.Join(tmpDir, "http-chart")
		result, err := PullChart(context.Background(), httpServer.URL, "test-chart", dstDir)
		if err != nil {
			t.Fatalf("pulling chart from HTTP server: %v", err)
		}
		if result.Chart.Metadata.Name != "test-chart" {
			t.Errorf("expected chart name 'test-chart', got %q", result.Chart.Metadata.Name)
		}
		if result.Version != "0.1.0" {
			t.Errorf("expected chart version '0.1.0', got %q", result.Version)
		}
	})

	// Start a local OCI registry server to serve the chart tarball
	ociServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/v2/library/test-chart/tags/list"):
			w.Header().Set("Content-Type", "application/vnd.oci.image.manifest.v1+json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"name":"library/test-chart","tags":["0.1.0", "latest"]}`)
		case strings.HasPrefix(r.URL.Path, "/v2/library/test-chart/manifests/0.1.0"):
			w.Header().Set("Content-Type", "application/vnd.oci.image.manifest.v1+json")
			http.ServeFile(w, r, filepath.Join(testdataDir, "manifest.json"))
		case strings.HasPrefix(r.URL.Path, "/v2/library/test-chart/manifests/sha256:"):
			w.Header().Set("Content-Type", "application/vnd.cncf.helm.config.v1+json")
			http.ServeFile(w, r, filepath.Join(testdataDir, "manifest.json"))
		case strings.HasPrefix(r.URL.Path, "/v2/library/test-chart/blobs/sha256:0325cff29a08050735ece813f6f37a077bc048484ca3516df0b2cb935114694b"):
			w.Header().Set("Content-Type", "application/vnd.cncf.helm.config.v1+json")
			http.ServeFile(w, r, filepath.Join(testdataDir, "chart.json"))
		case strings.HasPrefix(r.URL.Path, "/v2/library/test-chart/blobs/sha256:864b21413ad46ca48d418f3a1d2235426be79f6d273e9222abbb5012df91424e"):
			w.Header().Set("Docker-Content-Digest", "sha256:864b21413ad46ca48d418f3a1d2235426be79f6d273e9222abbb5012df91424e")
			w.Header().Set("Content-Type", "application/vnd.cncf.helm.chart.content.v1.tar+gzip")
			w.Header().Set("Content-Length", "4169")
			http.ServeFile(w, r, filepath.Join(testdataDir, "test-chart-0.1.0.tgz"))
		case strings.HasPrefix(r.URL.Path, "/library/test-chart"):
			http.ServeFile(w, r, filepath.Join(testdataDir, "test-chart-0.1.0.tgz"))
		}
	}))
	defer ociServer.Close()

	// Test pulling the chart from the OCI registry server
	t.Run("OCI registry", func(t *testing.T) {
		dstDir := filepath.Join(tmpDir, "oci-chart")
		registryClient, err := registry.NewClient(registry.ClientOptPlainHTTP(), registry.ClientOptDebug(true))
		if err != nil {
			t.Fatalf("creating registry client: %v", err)
		}
		ociServer.URL = strings.Replace(ociServer.URL, "http://", "oci://", 1)
		ociServer.URL = fmt.Sprintf("%s/library/test-chart", ociServer.URL)
		result, err := PullChart(context.Background(), ociServer.URL, "test-chart", dstDir, WithRegistryClient(registryClient))
		if err != nil {
			t.Fatalf("pulling chart from OCI registry: %v", err)
		}
		if result.Chart.Metadata.Name != "test-chart" {
			t.Errorf("expected chart name 'test-chart', got %q", result.Chart.Metadata.Name)
		}
		if result.Version != "0.1.0" {
			t.Errorf("expected chart version '0.1.0', got %q", result.Version)
		}
	})

	// Test pulling the chart with a specific version
	t.Run("WithVersion option", func(t *testing.T) {
		dstDir := filepath.Join(tmpDir, "version-chart")
		result, err := PullChart(context.Background(), httpServer.URL, "test-chart", dstDir, WithVersion("0.1.0"))
		if err != nil {
			t.Fatalf("pulling chart with version: %v", err)
		}
		if result.Chart.Metadata.Name != "test-chart" {
			t.Errorf("expected chart name 'test-chart', got %q", result.Chart.Metadata.Name)
		}
		if result.Version != "0.1.0" {
			t.Errorf("expected chart version '0.1.0', got %q", result.Version)
		}
	})
}
