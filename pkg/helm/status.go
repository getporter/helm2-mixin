package helm

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/deislabs/porter/pkg/printer"
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
func (m *Mixin) Status(opts printer.PrintOptions) error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	var step StatusStep
	err = yaml.Unmarshal(payload, &step)
	if err != nil {
		return err
	}

	format := ""
	switch opts.Format {
	case printer.FormatPlaintext:
		// do nothing, as default output is plaintext
	case printer.FormatYaml:
		format = `-o yaml`
	case printer.FormatJson:
		format = `-o json`
	default:
		return fmt.Errorf("invalid format: %s", opts.Format)
	}

	for _, release := range step.Arguments.Releases {
		cmd := m.NewCommand("helm", "status", strings.TrimSpace(fmt.Sprintf(`%s %s`, release, format)))

		cmd.Stdout = m.Out
		cmd.Stderr = m.Err

		prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
		fmt.Fprintln(m.Out, prettyCmd)

		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("could not execute command, %s: %s", prettyCmd, err)
		}
		err = cmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
