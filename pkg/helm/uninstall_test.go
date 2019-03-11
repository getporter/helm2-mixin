package helm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

type UninstallTest struct {
	expectedCommand string
	uninstallStep   UninstallStep
}

func TestMixin_UnmarshalUninstallStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/uninstall-input.yaml")
	require.NoError(t, err)

	var action UninstallAction
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	require.Len(t, action.Steps, 1)
	step := action.Steps[0]

	assert.Equal(t, "Uninstall MySQL", step.Description)
	assert.Equal(t, []string{"porter-ci-mysql"}, step.Releases)
	assert.True(t, step.Purge)
}

func TestMixin_Uninstall(t *testing.T) {
	releases := []string{
		"foo",
		"bar",
	}

	uninstallTests := []UninstallTest{
		{
			expectedCommand: `helm delete foo bar`,
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step:     Step{Description: "Uninstall Foo"},
					Releases: releases,
				},
			},
		},
		{
			expectedCommand: `helm delete --purge foo bar`,
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step:     Step{Description: "Uninstall Foo"},
					Purge:    true,
					Releases: releases,
				},
			},
		},
	}

	defer os.Unsetenv(test.ExpectedCommandEnv)
	for _, uninstallTest := range uninstallTests {
		t.Run(uninstallTest.expectedCommand, func(t *testing.T) {
			os.Setenv(test.ExpectedCommandEnv, uninstallTest.expectedCommand)

			action := UninstallAction{Steps: []UninstallStep{uninstallTest.uninstallStep}}
			b, _ := yaml.Marshal(action)

			x := string(b)
			fmt.Println(x)
			h := NewTestMixin(t)
			h.In = bytes.NewReader(b)

			err := h.Uninstall()

			require.NoError(t, err)
		})
	}
}
