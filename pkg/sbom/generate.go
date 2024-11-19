package sbom

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/sbom"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

// GenerateSBOMsFromManifest generates SBOMs from the given manifest
func GenerateSBOMsFromManifest(ctx context.Context, manifest io.Reader) ([]*sbom.SBOM, error) {
	imageNames, err := extractImageNamesFromManifest(ctx, manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to extract image names from manifest: %w", err)
	}

	var sboms []*sbom.SBOM
	for _, imageName := range imageNames {
		cfg := syft.DefaultCreateSBOMConfig().WithParallelism(10)
		src, err := syft.GetSource(
			ctx,
			imageName,
			syft.DefaultGetSourceConfig(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get source for image %s: %w", imageName, err)
		}
		s, err := syft.CreateSBOM(ctx, src, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create SBOM for image %s: %w", imageName, err)
		}
		sboms = append(sboms, s)
	}

	return sboms, nil
}

// GenerateSBOMsFromPath generates SBOMs from all manifests in the given path
func GenerateSBOMsFromPath(ctx context.Context, path string) ([]*sbom.SBOM, error) {
	logger := slog.Default()
	manifests, err := filepath.Glob(filepath.Join(path, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to read manifests: %w", err)
	}
	var sboms []*sbom.SBOM
	for _, manifest := range manifests {
		f, err := os.Open(manifest)
		defer func() {
			err := f.Close()
			if err != nil {
				logger.ErrorContext(ctx, "failed to close manifest file", slog.String("file", manifest), slog.Any("error", err))
			}
		}()
		if err != nil {
			return nil, fmt.Errorf("failed to open manifest file %s: %w", manifest, err)
		}
		sbom, err := GenerateSBOMsFromManifest(ctx, f)
		if err != nil {
			return nil, fmt.Errorf("failed to generate SBOM from manifest file %s: %w", manifest, err)
		}
		sboms = append(sboms, sbom...)
	}
	return sboms, nil
}

func extractImageNamesFromManifest(ctx context.Context, manifest io.Reader) ([]string, error) {
	logger := slog.Default()
	imageNames := []string{}
	var m []byte
	_, err := manifest.Read(m)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var obj unstructured.Unstructured
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, _, err = dec.Decode(m, nil, &obj)
	if err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %v\n", err)
	}

	containers, err := getContainers(&obj)
	if err != nil {
		return nil, fmt.Errorf("failed to extract containers: %v\n", err)
	}

	for _, container := range containers {
		containerMap, ok := container.(map[string]interface{})
		if !ok {
			logger.DebugContext(ctx, "failed to convert container to map[string]interface{}")
			continue
		}
		imageName, found, err := unstructured.NestedString(containerMap, "image")
		if err != nil {
			logger.DebugContext(ctx, "failed to extract image name from container", slog.Any("error", err))
			continue
		}
		if found {
			imageNames = append(imageNames, imageName)
		}
	}

	return imageNames, nil
}

func getContainers(obj *unstructured.Unstructured) ([]interface{}, error) {
	kind := obj.GetKind()
	var (
		containers []interface{}
		err        error
	)
	if kind == "Pod" {
		containers, _, err = unstructured.NestedSlice(obj.Object, "spec", "containers")
	} else {
		containers, _, err = unstructured.NestedSlice(obj.Object, "spec", "template",
			"spec", "containers")
	}
	if err != nil {
		return nil, err
	}
	return containers, nil
}
