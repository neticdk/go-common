package sbom

import (
	"bytes"
	"testing"

	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		format  Format
		wantErr bool
	}{
		{
			name:    "SPDX JSON",
			format:  FormatSPDXJSON,
			wantErr: false,
		},
		{
			name:    "SPDX Tag-Value",
			format:  FormatSPDXTagValue,
			wantErr: false,
		},
		{
			name:    "CycloneDX JSON",
			format:  FormatCycloneDXJSON,
			wantErr: false,
		},
		{
			name:    "CycloneDX XML",
			format:  FormatCycloneDXXML,
			wantErr: false,
		},
		{
			name:    "Unsupported Format",
			format:  Format(999),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			pkgs := []pkg.Package{
				{
					Name:    "example-package-1",
					Version: "1.0.0",
					Type:    "deb",
				},
				{
					Name:    "example-package-2",
					Version: "2.0.0",
					Type:    "rpm",
				},
				{
					Name:    "example-package-3",
					Version: "3.0.0",
					Type:    "npm",
				},
			}
			testSBOM := sbom.SBOM{
				Artifacts: sbom.Artifacts{
					Packages: pkg.NewCollection(pkgs...),
				},
				Source: source.Description{
					Name:    "example-source",
					Version: "1.0.0",
					Metadata: map[string]interface{}{
						"image": "example-image",
					},
				},
				Descriptor: sbom.Descriptor{
					Name:    "example-sbom",
					Version: "1.0.0",
				},
			}

			err := Encode(&buf, testSBOM, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("Encoded SBOM: %s", buf.String())

			if !tt.wantErr && buf.Len() == 0 {
				t.Errorf("Encode() = empty buffer, want non-empty buffer")
			}
		})
	}
}
