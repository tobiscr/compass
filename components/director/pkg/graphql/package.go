package graphql

type Package struct {
	ID                             string      `json:"id"`
	Name                           string      `json:"name"`
	Description                    *string     `json:"description"`
	InstanceAuthRequestInputSchema *JSONSchema `json:"InstanceAuthRequestInputSchema"`
	// When defined, all Auth requests fallback to defaultAuth.
	DefaultInstanceAuth *Auth `json:"defaultInstanceAuth"`
	AppID               string
}

type PackageExt struct {
	Package
	APIDefinitions           APIDefinitionPageExt      `json:"apiDefinitions"`
	EventDefinitions         EventAPIDefinitionPageExt `json:"eventDefinitions"`
	Documents                DocumentPageExt           `json:"documents"`
	APIDefinition            APIDefinitionExt          `json:"apiDefinition"`
	EventDefinition          EventDefinition           `json:"eventDefinition"`
	Document                 Document                  `json:"document"`
	InstanceAuth             *PackageInstanceAuth      `json:"instanceAuth"`
	InstanceAuths            []*PackageInstanceAuth    `json:"instanceAuths"`
	APIDefinitionsNoPaging   []*APIDefinitionExt       `json:"apiDefinitionsNoPaging"`
	EventDefinitionsNoPaging []*EventAPIDefinitionExt  `json:"eventDefinitionsNoPaging"`
}

type PackagePageExt struct {
	PackagePage
	Data []*PackageExt `json:"data"`
}
