package helm

import (
	"fmt"
	"strings"

	"get.porter.sh/porter/pkg/exec/builder"
	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// clientVersionConstraint represents the semver constraint for the Helm client version
// Currently, this mixin only supports Helm clients versioned v2.x.x
const clientVersionConstraint string = "^v2.x"

// These values may be referenced elsewhere (init.go), hence consts
const helmArchiveTmpl string = "helm-%s-linux-amd64.tar.gz"
const helmDownloadURLTmpl string = "https://get.helm.sh/%s"

const getHelm string = `RUN apt-get update && \
 apt-get install -y curl && \
 curl -o helm.tgz %s && \
 tar -xzf helm.tgz && \
 mv linux-amd64/helm /usr/local/bin && \
 rm helm.tgz
RUN helm init --client-only
`

// kubectl may be necessary; for example, to set up RBAC for Helm's Tiller component if needed
const kubeVersion string = "v1.15.3"
const getKubectl string = `RUN apt-get update && \
 apt-get install -y apt-transport-https curl && \
 curl -o kubectl https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/kubectl && \
 mv kubectl /usr/local/bin && \
 chmod a+x /usr/local/bin/kubectl`

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the helm mixin in porter.yaml
// mixins:
// - helm:
//	  repositories:
//	    stable:
//		  url: "https://charts.helm.sh/stable"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
	Repositories  map[string]Repository
}

type Repository struct {
	URL string `yaml:"url,omitempty"`
}

func (m *Mixin) Build() error {

	// Create new Builder.
	var input BuildInput
	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	suppliedClientVersion := input.Config.ClientVersion
	if suppliedClientVersion != "" {
		ok, err := validate(suppliedClientVersion, clientVersionConstraint)
		if err != nil {
			return err
		}
		if !ok {
			return errors.Errorf("supplied clientVersion %q does not meet semver constraint %q",
				suppliedClientVersion, clientVersionConstraint)
		}
		m.HelmClientVersion = suppliedClientVersion
	}

	var helmArchiveVersion = fmt.Sprintf(helmArchiveTmpl, m.HelmClientVersion)
	var helmDownloadURL = fmt.Sprintf(helmDownloadURLTmpl, helmArchiveVersion)

	// Define helm
	fmt.Fprintf(m.Out, getHelm, helmDownloadURL)

	// Define kubectl
	fmt.Fprintf(m.Out, getKubectl, kubeVersion)

	// Go through repositories if defined
	if len(input.Config.Repositories) > 0 {
		// Add the repositories
		for name, repo := range input.Config.Repositories {
			url := repo.URL
			repositoryCommand, err := getRepositoryCommand(name, url)
			if err != nil && m.Debug {
				fmt.Fprintf(m.Err, "DEBUG: addition of repository failed: %s\n", err.Error())
			} else {
				fmt.Fprintf(m.Out, strings.Join(repositoryCommand, " "))
			}
		}
		// Make sure we update the helm repositories
		// So we don't have to do it at runtime
		fmt.Fprintf(m.Out, "\nRUN helm repo update")
	}

	return nil
}

func getRepositoryCommand(name, url string) (repositoryCommand []string, err error) {

	var commandBuilder []string

	if url == "" {
		return commandBuilder, fmt.Errorf("repository url must be supplied")
	}

	commandBuilder = append(commandBuilder, "\nRUN", "helm", "repo", "add", name, url)

	return commandBuilder, nil
}

// validate validates that the supplied clientVersion meets the supplied semver constraint
func validate(clientVersion, constraint string) (bool, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, errors.Wrapf(err, "unable to parse version constraint %q", constraint)
	}

	v, err := semver.NewVersion(clientVersion)
	if err != nil {
		return false, errors.Wrapf(err, "supplied client version %q cannot be parsed as semver", clientVersion)
	}

	return c.Check(v), nil
}
