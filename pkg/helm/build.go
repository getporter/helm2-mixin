package helm

import (
	"fmt"
	"strings"

	"get.porter.sh/porter/pkg/exec/builder"
	yaml "gopkg.in/yaml.v2"
)

// These values may be referenced elsewhere (init.go), hence consts
const helmClientVersion string = "v2.15.2"
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

var helmArchiveVersion = fmt.Sprintf(helmArchiveTmpl, helmClientVersion)
var helmDownloadURL = fmt.Sprintf(helmDownloadURLTmpl, helmArchiveVersion)

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
//		  url: "https://kubernetes-charts.storage.googleapis.com"
//		  cafile: "path/to/cafile"
//		  certfile: "path/to/certfile"
//		  keyfile: "path/to/keyfile"
//		  username: "username"
//		  password: "password"
type MixinConfig struct {
	Repositories map[string]Repository
}

type Repository struct {
	URL      string `yaml:"url,omitempty"`
	Cafile   string `yaml:"cafile,omitempty"`
	Certfile string `yaml:"certfile,omitempty"`
	Keyfile  string `yaml:"keyfile,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
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

	// Define helm
	fmt.Fprintf(m.Out, getHelm, helmDownloadURL)

	// Define kubectl
	fmt.Fprintf(m.Out, getKubectl, kubeVersion)

	// Go through repositories
	for name, repo := range input.Config.Repositories {

		commandValue, err := GetAddRepositoryCommand(name, repo.URL, repo.Cafile, repo.Certfile, repo.Keyfile, repo.Username, repo.Password)
		if err != nil && m.Debug {
			fmt.Fprintf(m.Err, err.Error())
		} else {
			fmt.Fprintf(m.Out, strings.Join(commandValue, " "))
		}
	}

	return nil
}

func GetAddRepositoryCommand(name, url, cafile, certfile, keyfile, username, password string) (commandValue []string, err error) {

	var commandBuilder []string

	if url == "" {
		return commandBuilder, fmt.Errorf("\nRepository url must be supplied")
	}

	commandBuilder = append(commandBuilder, "\nRUN", "helm", "repo", "add", name, url)

	if certfile != "" && keyfile != "" {
		commandBuilder = append(commandBuilder, "--cert-file", certfile, "--key-file", keyfile)
	}
	if cafile != "" {
		commandBuilder = append(commandBuilder, "--ca-file", cafile)
	}
	if username != "" && password != "" {
		commandBuilder = append(commandBuilder, "--username", username, "--password", password)
	}

	return commandBuilder, nil
}
