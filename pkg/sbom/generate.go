package sbom

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/sbom"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

// GenerateSBOMFromManifests generates SBOMs from the given manifest files
func GenerateSBOMFromManifests(manifestFiles []string) ([]*sbom.SBOM, error) {
	ctx := context.Background()
	imageNames := extractImageNamesFromManifests(manifestFiles)
	fmt.Printf("Extracted image names: %v\n", imageNames)

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

// GenerateSBOMFromPath generates SBOMs from all manifests in the given path
func GenerateSBOMFromPath(path string) ([]*sbom.SBOM, error) {
	manifests, err := filepath.Glob(filepath.Join(path, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to read manifests: %w", err)
	}
	return GenerateSBOMFromManifests(manifests)
}

func extractImageNamesFromManifests(manifestFiles []string) []string {
	imageNames := []string{}
	for _, manifestFile := range manifestFiles {
		file, err := os.ReadFile(manifestFile)
		if err != nil {
			fmt.Printf("failed to read manifest file %s: %v\n", manifestFile, err)
			continue
		}

		var obj unstructured.Unstructured
		dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		_, _, err = dec.Decode(file, nil, &obj)
		if err != nil {
			fmt.Printf("failed to decode manifest file %s: %v\n", manifestFile, err)
			continue
		}

		containers, err := getContainers(&obj)
		if err != nil {
			fmt.Printf("failed to extract containers from manifest file %s: %v\n", manifestFile, err)
			continue
		}

		if len(containers) == 0 {
			continue
		}

		for _, container := range containers {
			containerMap, ok := container.(map[string]interface{})
			if !ok {
				fmt.Printf("failed to convert container to map[string]interface{}\n")
				continue
			}
			imageName, found, err := unstructured.NestedString(containerMap, "image")
			if err != nil {
				fmt.Printf("failed to extract image name from container: %v\n", err)
				continue
			}
			if found {
				imageNames = append(imageNames, imageName)
			}
		}
	}
	return imageNames
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
