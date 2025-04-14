package types

type SBOMCDX struct {
	BOMFormat    string       `json,yaml:"bomFormat"`
	SpecVersion  string       `json,yaml:"specVersion"`
	SerialNumber string       `json,yaml:"serialNumber"`
	Version      int          `json,yaml:"version"`
	Metadata     Metadata     `json,yaml:"metadata"`
	Components   []Component  `json,yaml:"components"`
	Dependencies []Dependency `json,yaml:"dependencies"`
}

type Metadata struct {
	Timestamp string    `json,yaml:"timestamp"`
	Tools     Tool      `json,yaml:"tools"`
	Component Component `json,yaml:"component"`
}

type Tool struct {
	Vendor     string           `json,yaml:"vendor"`
	Name       string           `json,yaml:"name"`
	Version    string           `json,yaml:"version"`
	Components []ToolsComponent `json,yaml:"components"`
}

type Component struct {
	SBOMRef     string     `json,yaml:"bom-ref"`
	Type        string     `json,yaml:"type"`
	Name        string     `json,yaml:"name"`
	Purl        string     `json,yaml:"purl"`
	Version     string     `json,yaml:"version"`
	Description string     `json,yaml:"description"`
	Licenses    []License  `json,yaml:"licenses"`
	Properties  []Property `json,yaml:"properties"`
}

type ToolsComponent struct {
	Type    string `json,yaml:"type"`
	Group   string `json,yaml:"group"`
	Name    string `json,yaml:"name"`
	Version string `json,yaml:"version"`
}

type License struct {
	ID   string `json,yaml:"id"`
	Name string `json,yaml:"name"`
	Text string `json,yaml:"text"`
}

type Property struct {
	Name  string `json,yaml:"name"`
	Value string `json,yaml:"value"`
}

type Dependency struct {
	Ref       string   `json,yaml:"ref"`
	DependsOn []string `json,yaml:"dependsOn"`
}
