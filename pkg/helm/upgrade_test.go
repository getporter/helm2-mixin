package helm

import (
	"bytes"
	"os"
	"testing"

	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestMixin_Upgrade(t *testing.T) {
	os.Setenv(test.ExpectedCommandEnv, `helm upgrade MYRELEASE MYCHART --namespace MYNAMESPACE --version 1.0.0 --replace --values /tmp/val1.yaml --values /tmp/val2.yaml --set baz=qux --set foo=bar`)
	defer os.Unsetenv(test.ExpectedCommandEnv)

	step := UpgradeStep{
		Arguments: UpgradeArguments{
			Namespace: "MYNAMESPACE",
			Name:      "MYRELEASE",
			Chart:     "MYCHART",
			Version:   "1.0.0",
			Replace:   true,
			Set: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			Values: []string{
				"/tmp/val1.yaml",
				"/tmp/val2.yaml",
			},
		},
	}
	b, _ := yaml.Marshal(step)

	h := NewTestMixin(t)
	h.In = bytes.NewReader(b)

	err := h.Upgrade()

	require.NoError(t, err)
}
