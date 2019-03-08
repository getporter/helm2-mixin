package helm

type Step struct {
	Description string       `yaml:"description"`
	Outputs     []HelmOutput `yaml:"outputs,omitempty"`
}

type HelmOutput struct {
	Name   string `yaml:"name"`
	Secret string `yaml:"secret"`
	Key    string `yaml:"key"`
}
