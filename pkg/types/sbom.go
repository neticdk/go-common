package types

// This file contains the definition of the CycloneDX Software Bill of Materials (SBOM) format.
// The CycloneDX format is a lightweight SBOM standard designed for use in application security
// contexts and supply chain component analysis. Here a simple edition of that is defined.

// The SBOMCDX struct represents the CycloneDX Software Bill of Materials (SBOM) format.
type SBOMCDX struct {
	BOMFormat    string       `json:"bomFormat" yaml:"bomFormat"`
	SpecVersion  string       `json:"specVersion" yaml:"specVersion"`
	SerialNumber string       `json:"serialNumber" yaml:"serialNumber"`
	Version      int          `json:"version" yaml:"version"`
	Metadata     Metadata     `json:"metadata" yaml:"metadata"`
	Components   []Component  `json:"components" yaml:"components"`
	Dependencies []Dependency `json:"dependencies" yaml:"dependencies"`
}

// Metadata represents the metadata of the SBOM
type Metadata struct {
	Timestamp string    `json:"timestamp" yaml:"timestamp"`
	Tools     Tool      `json:"tools" yaml:"tools"`
	Component Component `json:"component" yaml:"component"`
}

// Tool represents the tool information of the SBOM
type Tool struct {
	Vendor     string           `json:"vendor" yaml:"vendor"`
	Name       string           `json:"name" yaml:"name"`
	Version    string           `json:"version" yaml:"version"`
	Components []ToolsComponent `json:"components" yaml:"components"`
}

// Component represents a component information of the SBOM
type Component struct {
	SBOMRef     string     `json:"bom-ref" yaml:"bom-ref"`
	Type        string     `json:"type" yaml:"type"`
	Name        string     `json:"name" yaml:"name"`
	Purl        string     `json:"purl" yaml:"purl"`
	Version     string     `json:"version" yaml:"version"`
	Description string     `json:"description" yaml:"description"`
	Licenses    []License  `json:"licenses" yaml:"licenses"`
	Properties  []Property `json:"properties" yaml:"properties"`
}

// ToolsComponent represents a component part of the SBOM tool
type ToolsComponent struct {
	Type    string `json:"type" yaml:"type"`
	Group   string `json:"group" yaml:"group"`
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

// License representartion
type License struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Text string `json:"text" yaml:"text"`
}

// Property represents a generic name value construct
type Property struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// Dependency represents  a dependency in SBOM
type Dependency struct {
	Ref       string   `json:"ref" yaml:"ref"`
	DependsOn []string `json:"dependsOn" yaml:"dependsOn"`
}
