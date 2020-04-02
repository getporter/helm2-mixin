//go:generate packr2

package helm

import (
	"bufio"
	"io/ioutil"
	"strings"

	"get.porter.sh/mixin/helm/pkg/kubernetes"
	"get.porter.sh/porter/pkg/context"
	"github.com/ghodss/yaml" // We are not using go-yaml because of serialization problems with jsonschema, don't use this library elsewhere
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	k8s "k8s.io/client-go/kubernetes"
)

const defaultHelmClientVersion string = "v2.15.2"

// Helm is the logic behind the helm mixin
type Mixin struct {
	*context.Context
	schema        *packr.Box
	ClientFactory kubernetes.ClientFactory
	TillerIniter
	HelmClientVersion string
}

// New helm mixin client, initialized with useful defaults.
func New() *Mixin {
	return &Mixin{
		schema:            packr.New("schema", "./schema"),
		Context:           context.New(),
		ClientFactory:     kubernetes.New(),
		TillerIniter:      RealTillerIniter{},
		HelmClientVersion: defaultHelmClientVersion,
	}
}

func (m *Mixin) getPayloadData() ([]byte, error) {
	reader := bufio.NewReader(m.In)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "could not read the payload from STDIN")
	}
	return data, nil
}

func (m *Mixin) ValidatePayload(b []byte) error {
	// Load the step as a go dump
	s := make(map[string]interface{})
	err := yaml.Unmarshal(b, &s)
	if err != nil {
		return errors.Wrap(err, "could not marshal payload as yaml")
	}
	manifestLoader := gojsonschema.NewGoLoader(s)

	// Load the step schema
	schema, err := m.GetSchema()
	if err != nil {
		return err
	}
	schemaLoader := gojsonschema.NewStringLoader(schema)

	validator, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return errors.Wrap(err, "unable to compile the mixin step schema")
	}

	// Validate the manifest against the schema
	result, err := validator.Validate(manifestLoader)
	if err != nil {
		return errors.Wrap(err, "unable to validate the mixin step schema")
	}
	if !result.Valid() {
		errs := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			errs = append(errs, err.String())
		}
		return errors.New(strings.Join(errs, "\n\t* "))
	}

	return nil
}

func (m *Mixin) getKubernetesClient(kubeconfig string) (k8s.Interface, error) {
	return m.ClientFactory.GetClient(kubeconfig)
}

func (m *Mixin) getHelmClientVersion() string {
	return m.HelmClientVersion
}

func (m *Mixin) setHelmClientVersion(version string) {
	if version != "" {
		m.HelmClientVersion = version
	}
}
