package helm

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMixin_PrintSchema(t *testing.T) {
	m := NewTestMixin(t)

	err := m.PrintSchema()
	require.NoError(t, err)

	gotSchema := m.TestContext.GetOutput()

	wantSchema, err := ioutil.ReadFile("testdata/schema.json")
	require.NoError(t, err)

	assert.Equal(t, string(wantSchema), gotSchema)
}

func TestMixin_ValidatePayload(t *testing.T) {
	testcases := []struct {
		name  string
		step  string
		pass  bool
		error string
	}{
		{"install", "testdata/install-input.yaml", true, ""},
		{"status", "testdata/status-input.yaml", true, ""},
		{"upgrade", "testdata/upgrade-input.yaml", true, ""},
		{"uninstall", "testdata/status-input.yaml", true, ""},
		{"install.missing-desc", "testdata/bad-install-input.missing-desc.yaml", false, "install.0.helm.description: String length must be greater than or equal to 1"},
		{"status.missing-releases", "testdata/bad-status-input.missing-releases.yaml", false, "status.0.helm: releases is required"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewTestMixin(t)
			b, err := ioutil.ReadFile(tc.step)
			require.NoError(t, err)

			err = m.ValidatePayload(b)
			if tc.pass {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.error)
			}
		})
	}
}
