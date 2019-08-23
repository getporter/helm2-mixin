package helm

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type InstallAction struct {
	Steps []InstallStep `yaml:"install"`
}

type InstallStep struct {
	InstallArguments `yaml:"helm"`
}

type InstallArguments struct {
	Step `yaml:",inline"`

	Namespace string            `yaml:"namespace"`
	Name      string            `yaml:"name"`
	Chart     string            `yaml:"chart"`
	Version   string            `yaml:"version"`
	Replace   bool              `yaml:"replace"`
	Set       map[string]string `yaml:"set"`
	Values    []string          `yaml:"values"`
	Devel     bool              `yaml:"devel`
	Wait      bool              `yaml:"wait"`
}

func (m *Mixin) Install() error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	kubeClient, err := m.getKubernetesClient("/root/.kube/config")
	if err != nil {
		return errors.Wrap(err, "couldn't get kubernetes client")
	}

	var action InstallAction
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

	cmd := m.NewCommand("helm", "install", "--name", step.Name, step.Chart)

	if step.Namespace != "" {
		cmd.Args = append(cmd.Args, "--namespace", step.Namespace)
	}

	if step.Version != "" {
		cmd.Args = append(cmd.Args, "--version", step.Version)
	}

	if step.Replace {
		cmd.Args = append(cmd.Args, "--replace")
	}

	if step.Wait {
		cmd.Args = append(cmd.Args, "--wait")
	}

	if step.Devel {
		cmd.Args = append(cmd.Args, "--devel")
	}

	for _, v := range step.Values {
		cmd.Args = append(cmd.Args, "--values", v)
	}

	// sort the set consistently
	setKeys := make([]string, 0, len(step.Set))
	for k := range step.Set {
		setKeys = append(setKeys, k)
	}
	sort.Strings(setKeys)

	for _, k := range setKeys {
		cmd.Args = append(cmd.Args, "--set", fmt.Sprintf("%s=%s", k, step.Set[k]))
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

	for _, output := range step.Outputs {
		val, err := getSecret(kubeClient, step.Namespace, output.Secret, output.Key)
		if err != nil {
			return err
		}

		err = m.Context.WriteMixinOutputToFile(output.Name, val)
		if err != nil {
			return errors.Wrapf(err, "unable to write output '%s'", output.Name)
		}
	}
	return nil
}
