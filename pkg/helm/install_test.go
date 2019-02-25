package helm

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

type InstallTest struct {
	expectedCommand string
	installStep     InstallStep
}

// sad hack: not sure how to make a common test main for all my subpackages
func TestMain(m *testing.M) {
	test.TestMainWithMockedCommandHandlers(m)
}

func TestMixin_Install(t *testing.T) {
	namespace := "MYNAMESPACE"
	name := "MYRELEASE"
	chart := "MYCHART"
	version := "1.0.0"
	setArgs := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}
	values := []string{
		"/tmp/val1.yaml",
		"/tmp/val2.yaml",
	}

	baseInstall := fmt.Sprintf(`helm install --name %s %s --namespace %s --version %s`, name, chart, namespace, version)
	baseValues := `--values /tmp/val1.yaml --values /tmp/val2.yaml`
	baseSetArgs := `--set baz=qux --set foo=bar`

	installTests := []InstallTest{
		{
			expectedCommand: fmt.Sprintf(`%s %s %s`, baseInstall, baseValues, baseSetArgs),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
					Namespace: namespace,
					Name:      name,
					Chart:     chart,
					Version:   version,
					Set:       setArgs,
					Values:    values,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf(`%s %s %s %s`, baseInstall, `--replace`, baseValues, baseSetArgs),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
					Namespace: namespace,
					Name:      name,
					Chart:     chart,
					Version:   version,
					Set:       setArgs,
					Values:    values,
					Replace:   true,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf(`%s %s %s %s`, baseInstall, `--wait`, baseValues, baseSetArgs),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
					Namespace: namespace,
					Name:      name,
					Chart:     chart,
					Version:   version,
					Set:       setArgs,
					Values:    values,
					Wait:      true,
				},
			},
		},
	}

	for _, installTest := range installTests {
		os.Setenv(test.ExpectedCommandEnv, installTest.expectedCommand)
		defer os.Unsetenv(test.ExpectedCommandEnv)

		b, _ := yaml.Marshal(installTest.installStep)

		h := NewTestMixin(t)
		h.In = bytes.NewReader(b)

		err := h.Install()

		require.NoError(t, err)
	}
}
