package helm

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type UninstallAction struct {
	Steps []UninstallStep `yaml:"uninstall"`
}

// UninstallStep represents the structure of an Uninstall action
type UninstallStep struct {
	UninstallArguments `yaml:"helm"`
}

// UninstallArguments are the arguments available for the Uninstall action
type UninstallArguments struct {
	Step `yaml:",inline"`

	Releases []string `yaml:"releases"`
	Purge    bool     `yaml:"purge"`
}

// Uninstall deletes a provided set of Helm releases, supplying optional flags/params
func (m *Mixin) Uninstall() error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	var action UninstallAction
	err = yaml.Unmarshal(payload, &action)
	if err != nil {
		return err
	}
	if len(action.Steps) != 1 {
		return errors.Errorf("expected a single step, but got %d", len(action.Steps))
	}
	step := action.Steps[0]

	err = m.Init()
	if err != nil {
		return err
	}

	// Delete each release one at a time, because helm stops on first error
	// This gives us more fine-grained error recovery and handling
	var result error
	for _, release := range step.Releases {
		err = m.delete(release, step.Purge)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

func (m *Mixin) delete(release string, purge bool) error {
	cmd := m.NewCommand("helm", "delete")

	if purge {
		cmd.Args = append(cmd.Args, "--purge")
	}

	cmd.Args = append(cmd.Args, release)

	output := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(m.Out, output)
	cmd.Stderr = io.MultiWriter(m.Err, output)

	prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
	fmt.Fprintln(m.Out, prettyCmd)

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("could not execute command, %s: %s", prettyCmd, err)
	}
	err = cmd.Wait()
	if err != nil {
		// Gracefully handle the error being a release not found
		if strings.Contains(output.String(), fmt.Sprintf(`release: %q not found`, release)) {
			return nil
		}
		return err
	}

	return nil
}
