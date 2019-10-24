package helm

import (
	"github.com/deislabs/porter/pkg/exec/builder"
)

var _ builder.ExecutableAction = Action{}

type Action struct {
	Steps []ExecuteSteps // using UnmarshalYAML so that we don't need a custom type per action
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - helm: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var steps []ExecuteSteps
	results, err := builder.UnmarshalAction(unmarshal, &steps)
	if err != nil {
		return err
	}

	for _, result := range results {
		step := result.(*[]ExecuteSteps)
		a.Steps = append(a.Steps, *step...)
	}
	return nil
}

func (a Action) GetSteps() []builder.ExecutableStep {
	steps := make([]builder.ExecutableStep, len(a.Steps))
	for i := range a.Steps {
		steps[i] = a.Steps[i]
	}

	return steps
}

var _ builder.ExecutableStep = ExecuteStep{}

type ExecuteSteps struct {
	ExecuteStep `yaml:"helm"`
}

type ExecuteStep struct {
	Step      `yaml:",inline"`
	Namespace string        `yaml:"namespace,omitempty"`
	Arguments []string      `yaml:"arguments,omitempty"`
	Flags     builder.Flags `yaml:"flags,omitempty"`
}

func (s ExecuteStep) GetCommand() string {
	return "helm"
}

func (s ExecuteStep) GetArguments() []string {
	return s.Arguments
}

func (s ExecuteStep) GetFlags() builder.Flags {
	return s.Flags
}
