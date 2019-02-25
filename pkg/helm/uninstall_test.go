package helm

import (
	"bytes"
	"os"
	"testing"

	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

type UninstallTest struct {
	expectedCommand string
	uninstallStep   UninstallStep
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
					Releases: releases,
				},
			},
		},
		{
			expectedCommand: `helm delete --purge foo bar`,
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
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

			b, _ := yaml.Marshal(uninstallTest.uninstallStep)

			h := NewTestMixin(t)
			h.In = bytes.NewReader(b)

			err := h.Uninstall()

			require.NoError(t, err)
		})
	}
}
