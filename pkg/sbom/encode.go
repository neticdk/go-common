package sbom

import (
	"fmt"
	"io"

	"github.com/anchore/syft/syft/format/cyclonedxjson"
	"github.com/anchore/syft/syft/format/cyclonedxxml"
	"github.com/anchore/syft/syft/format/spdxjson"
	"github.com/anchore/syft/syft/format/spdxtagvalue"
	syftSbom "github.com/anchore/syft/syft/sbom"
)

type Format int

const (
	FormatSPDXJSON Format = iota
	FormatSPDXTagValue
	FormatCycloneDXJSON
	FormatCycloneDXXML
)

type Encoder interface {
	// Encode writes the SBOM to the writer in the specified format
	Encode(io.Writer, syftSbom.SBOM) error
}

// Encode writes the SBOM to the writer in the specified format
func Encode(w io.Writer, sbom syftSbom.SBOM, f Format) error {
	var encoder Encoder
	switch f {
	case FormatSPDXJSON:
		enc, err := spdxjson.NewFormatEncoderWithConfig(spdxjson.DefaultEncoderConfig())
		if err != nil {
			return fmt.Errorf("creating encoder: %w", err)
		}
		encoder = enc
	case FormatSPDXTagValue:
		enc, err := spdxtagvalue.NewFormatEncoderWithConfig(spdxtagvalue.DefaultEncoderConfig())
		if err != nil {
			return fmt.Errorf("creating encoder: %w", err)
		}
		encoder = enc
	case FormatCycloneDXJSON:
		enc, err := cyclonedxjson.NewFormatEncoderWithConfig(cyclonedxjson.DefaultEncoderConfig())
		if err != nil {
			return fmt.Errorf("creating encoder: %w", err)
		}
		encoder = enc
	case FormatCycloneDXXML:
		enc, err := cyclonedxxml.NewFormatEncoderWithConfig(cyclonedxxml.DefaultEncoderConfig())
		if err != nil {
			return fmt.Errorf("creating encoder: %w", err)
		}
		encoder = enc
	default:
		return fmt.Errorf("unsupported format: %v", f)
	}

	if encoder == nil {
		return fmt.Errorf("getting encoder for format: %v", f)
	}

	return encoder.Encode(w, sbom)
}
