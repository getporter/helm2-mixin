package helm2

type Step struct {
	Description string       `yaml:"description"`
	Outputs     []HelmOutput `yaml:"outputs,omitempty"`
}

type HelmOutput struct {
	Name         string `yaml:"name"`
	Secret       string `yaml:"secret,omitempty"`
	Key          string `yaml:"key,omitempty"`
	ResourceType string `yaml:"resourceType,omitempty"`
	ResourceName string `yaml:"resourceName,omitempty"`
	Namespace    string `yaml:"namespace,omitempty"`
	JSONPath     string `yaml:"jsonPath,omitempty"`
}
