package helm

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/deislabs/porter/pkg/exec/builder"
	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestMixin_UnmarshalExecuteStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/execute-input.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	require.Len(t, action.Steps, 1)
	step := action.Steps[0]

	assert.Equal(t, "MySQL Status", step.Description)
	assert.Equal(t, []string{"status", "mysql"}, step.Arguments)
	wantFlags := builder.Flags{
		builder.Flag{
			Name:   "o",
			Values: []string{"yaml"},
		},
	}
	assert.Equal(t, wantFlags, step.Flags)
	wantOutputs := []HelmOutput{
		{
			Name:   "mysql-root-password",
			Secret: "porter-ci-mysql",
			Key:    "mysql-root-password",
		},
	}
	assert.Equal(t, wantOutputs, step.Outputs)
}

func TestMixin_Execute(t *testing.T) {
	defer os.Unsetenv(test.ExpectedCommandEnv)
	os.Setenv(test.ExpectedCommandEnv, "helm status mysql -o yaml")

	executeAction := Action{
		Steps: []ExecuteSteps{
			{
				ExecuteStep: ExecuteStep{
					Arguments: []string{
						"status",
						"mysql",
					},
					Flags: builder.Flags{
						{
							Name: "o",
							Values: []string{
								"yaml",
							},
						},
					},
				},
			},
		},
	}

	b, _ := yaml.Marshal(executeAction)

	h := NewTestMixin(t)
	h.In = bytes.NewReader(b)

	err := h.Execute()
	require.NoError(t, err)
}
