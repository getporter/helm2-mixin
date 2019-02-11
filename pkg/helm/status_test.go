package helm

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/deislabs/porter/pkg/printer"
	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type statusTest struct {
	format                printer.Format
	expectedCommandSuffix string
}

func TestMixin_Status(t *testing.T) {
	testCases := map[string]statusTest{
		"default": statusTest{
			format:                printer.FormatPlaintext,
			expectedCommandSuffix: "",
		},
		"json": statusTest{
			format:                printer.FormatJson,
			expectedCommandSuffix: "-o json",
		},
		"yaml": statusTest{
			format:                printer.FormatYaml,
			expectedCommandSuffix: "-o yaml",
		},
	}

	releases := []string{
		"foo",
		"bar",
	}

	for _, testCase := range testCases {
		for _, release := range releases {
			os.Setenv(test.ExpectedCommandEnv,
				strings.TrimSpace(fmt.Sprintf(`helm status %s %s`, release, testCase.expectedCommandSuffix)))
			defer os.Unsetenv(test.ExpectedCommandEnv)

			statusStep := StatusStep{
				Arguments: StatusArguments{
					Releases: []string{release},
				},
			}

			b, _ := yaml.Marshal(statusStep)

			h := NewTestMixin(t)
			h.In = bytes.NewReader(b)

			opts := printer.PrintOptions{testCase.format}
			err := h.Status(opts)

			require.NoError(t, err)
		}
	}
}
