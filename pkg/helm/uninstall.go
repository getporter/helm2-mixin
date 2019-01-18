package helm

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// UninstallStep represents the structure of an Uninstall action
type UninstallStep struct {
	Description string             `yaml:"description"`
	Arguments   UninstallArguments `yaml:"helm"`
}

// UninstallArguments are the arguments available for the Uninstall action
type UninstallArguments struct {
	Releases []string `yaml:"releases"`
	Purge    bool     `yaml:"purge"`
}

// Uninstall deletes a provided set of Helm releases, supplying optional flags/params
func (m *Mixin) Uninstall() error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	var step UninstallStep
	err = yaml.Unmarshal(payload, &step)
	if err != nil {
		return err
	}

	cmd := m.NewCommand("helm", "delete")

	if step.Arguments.Purge {
		cmd.Args = append(cmd.Args, "--purge")
	}

	for _, release := range step.Arguments.Releases {
		cmd.Args = append(cmd.Args, release)
	}

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

	return nil
}
