# SBOM Package

The `sbom` package provides utilities for generating and encoding Software Bill
of Materials (SBOM) from Kubernetes manifests. It leverages the Syft library to
create SBOMs in various formats such as SPDX JSON, SPDX Tag-Value, CycloneDX
JSON, and CycloneDX XML.

## Installation

To use this package, you need to install it using `go get`:

```sh
go get github.com/yourusername/go-common/pkg/sbom
```

## Usage

### Generating SBOMs

You can generate SBOMs from Kubernetes manifests or from a directory containing
manifest files.

#### Generate SBOMs from a Manifest

To generate SBOMs from a Kubernetes manifest, use the
`GenerateSBOMsFromManifest` function:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yourusername/go-common/pkg/sbom"
)

func main() {
	ctx := context.Background()
	manifestFile, err := os.Open("path/to/manifest.yaml")
	if err != nil {
		fmt.Printf("Error opening manifest file: %v\n", err)
		return
	}
	defer manifestFile.Close()

	sboms, err := sbom.GenerateSBOMsFromManifest(ctx, manifestFile)
	if err != nil {
		fmt.Printf("Error generating SBOMs: %v\n", err)
		return
	}

	for _, s := range sboms {
		fmt.Printf("Generated SBOM: %+v\n", s)
	}
}
```

#### Generate SBOMs from a Path

To generate SBOMs from all manifest files in a directory, use the
`GenerateSBOMsFromPath` function:

```go
package main

import (
	"context"
	"fmt"

	"github.com/yourusername/go-common/pkg/sbom"
)

func main() {
	ctx := context.Background()
	path := "path/to/manifests"

	sboms, err := sbom.GenerateSBOMsFromPath(ctx, path)
	if err != nil {
		fmt.Printf("Error generating SBOMs: %v\n", err)
		return
	}

	for _, s := range sboms {
		fmt.Printf("Generated SBOM: %+v\n", s)
	}
}
```

### Encoding SBOMs

You can encode SBOMs into different formats using the `Encode` function.
Supported formats include SPDX JSON, SPDX Tag-Value, CycloneDX JSON, and
CycloneDX XML.

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yourusername/go-common/pkg/sbom"
)

func main() {
	ctx := context.Background()
	manifestFile, err := os.Open("path/to/manifest.yaml")
	if err != nil {
		fmt.Printf("Error opening manifest file: %v\n", err)
		return
	}
	defer manifestFile.Close()

	sboms, err := sbom.GenerateSBOMsFromManifest(ctx, manifestFile)
	if err != nil {
		fmt.Printf("Error generating SBOMs: %v\n", err)
		return
	}

	outputFile, err := os.Create("output.sbom.json")
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	err = sbom.Encode(outputFile, *sboms[0], sbom.FormatSPDXJSON)
	if err != nil {
		fmt.Printf("Error encoding SBOM: %v\n", err)
		return
	}

	fmt.Println("SBOM successfully encoded to output.sbom.json")
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License.

## Acknowledgements

This package uses the [Syft](https://github.com/anchore/syft) library for SBOM
generation and encoding. Special thanks to the Syft team for their excellent
work.
