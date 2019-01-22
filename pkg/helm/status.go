package helm

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// StatusStep represents the structure of an Status action
type StatusStep struct {
	Description string          `yaml:"description"`
	Arguments   StatusArguments `yaml:"helm"`
}

// StatusArguments are the arguments available for the Status action
type StatusArguments struct {
	Releases []string `yaml:"releases"`
}

// Status reports the status for a provided set of Helm releases
func (m *Mixin) Status() error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	var step StatusStep
	err = yaml.Unmarshal(payload, &step)
	if err != nil {
		return err
	}

	cmd := m.NewCommand("helm", "status")

	for _, release := range step.Arguments.Releases {
		statusCmd := cmd
		statusCmd.Args = append(statusCmd.Args, release)
		statusCmd.Stdout = m.Out
		statusCmd.Stderr = m.Err

		prettyCmd := fmt.Sprintf("%s %s", statusCmd.Path, strings.Join(statusCmd.Args, " "))
		fmt.Fprintln(m.Out, prettyCmd)

		err = statusCmd.Start()
		if err != nil {
			return fmt.Errorf("could not execute command, %s: %s", prettyCmd, err)
		}
		err = statusCmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
