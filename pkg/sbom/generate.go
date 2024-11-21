package sbom

import (
	"bufio"
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
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
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

	if path == "" {
		return make([]*sbom.SBOM, 0), nil
	}

	manifests, err := filepath.Glob(filepath.Join(path, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("reading manifests: %w", err)
	}
	var sboms []*sbom.SBOM
	for _, manifest := range manifests {
		f, err := os.Open(manifest)
		defer func() {
			err := f.Close()
			if err != nil {
				logger.ErrorContext(ctx, "closing manifest file", slog.String("file", manifest), slog.Any("error", err))
			}
		}()
		if err != nil {
			return nil, fmt.Errorf("opening manifest file %s: %v", manifest, err)
		}
		sbom, err := GenerateSBOMsFromManifest(ctx, f)
		if err != nil {
			return nil, fmt.Errorf("generating SBOM from manifest file %s: %v", manifest, err)
		}
		sboms = append(sboms, sbom...)
	}
	return sboms, nil
}

func extractImageNamesFromManifest(ctx context.Context, manifest io.Reader) ([]string, error) {
	logger := slog.Default()
	imageNames := []string{}

	multidocReader := utilyaml.NewYAMLReader(bufio.NewReader(manifest))
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	var containers []interface{}
	for {
		buf, err := multidocReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read manifest: %v\n", err)
		}

		var obj unstructured.Unstructured
		rto, _, err := dec.Decode(buf, nil, &obj)
		// happens with empty YAML documents
		if rto == nil {
			logger.DebugContext(ctx, "failed to decode manifest")
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to decode manifest: %v\n", err)
		}

		objContainers, err := getContainers(&obj)
		if err != nil {
			return nil, fmt.Errorf("failed to extract containers: %v\n", err)
		}
		containers = append(containers, objContainers...)
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
