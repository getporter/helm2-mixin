package helm

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type UpgradeTest struct {
	expectedCommand string
	upgradeStep     UpgradeStep
}

func TestMixin_Upgrade(t *testing.T) {
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

	baseUpgrade := fmt.Sprintf(`helm upgrade %s %s --namespace %s --version %s`, name, chart, namespace, version)
	baseValues := `--values /tmp/val1.yaml --values /tmp/val2.yaml`
	baseSetArgs := `--set baz=qux --set foo=bar`

	upgradeTests := []UpgradeTest{
		UpgradeTest{
			expectedCommand: fmt.Sprintf(`%s %s %s`, baseUpgrade, baseValues, baseSetArgs),
			upgradeStep: UpgradeStep{
				Arguments: UpgradeArguments{
					Namespace: namespace,
					Name:      name,
					Chart:     chart,
					Version:   version,
					Set:       setArgs,
					Values:    values,
				},
			},
		},
		UpgradeTest{
			expectedCommand: fmt.Sprintf(`%s %s %s %s`, baseUpgrade, `--reset-values`, baseValues, baseSetArgs),
			upgradeStep: UpgradeStep{
				Arguments: UpgradeArguments{
					Namespace:   namespace,
					Name:        name,
					Chart:       chart,
					Version:     version,
					Set:         setArgs,
					Values:      values,
					ResetValues: true,
				},
			},
		},
		UpgradeTest{
			expectedCommand: fmt.Sprintf(`%s %s %s %s`, baseUpgrade, `--reuse-values`, baseValues, baseSetArgs),
			upgradeStep: UpgradeStep{
				Arguments: UpgradeArguments{
					Namespace:   namespace,
					Name:        name,
					Chart:       chart,
					Version:     version,
					Set:         setArgs,
					Values:      values,
					ReuseValues: true,
				},
			},
		},
		UpgradeTest{
			expectedCommand: fmt.Sprintf(`%s %s %s %s`, baseUpgrade, `--wait`, baseValues, baseSetArgs),
			upgradeStep: UpgradeStep{
				Arguments: UpgradeArguments{
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

	for _, upgradeTest := range upgradeTests {
		os.Setenv(test.ExpectedCommandEnv, upgradeTest.expectedCommand)
		defer os.Unsetenv(test.ExpectedCommandEnv)

		b, err := yaml.Marshal(upgradeTest.upgradeStep)
		require.NoError(t, err)

		h := NewTestMixin(t)
		h.In = bytes.NewReader(b)

		err = h.Upgrade()

		require.NoError(t, err)
	}
}
