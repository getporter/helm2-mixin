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

func TestMixin_Status(t *testing.T) {
	releases := []string{
		"foo",
		"bar",
	}

	for _, release := range releases {
		os.Setenv(test.ExpectedCommandEnv, fmt.Sprintf(`helm status %s`, release))
		defer os.Unsetenv(test.ExpectedCommandEnv)

		statusStep := StatusStep{
			Arguments: StatusArguments{
				Releases: []string{release},
			},
		}

		b, _ := yaml.Marshal(statusStep)

		h := NewTestMixin(t)
		h.In = bytes.NewReader(b)

		err := h.Status()

		require.NoError(t, err)
	}
}
