package helm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMixin_Build(t *testing.T) {
	m := NewTestMixin(t)

	err := m.Build()
	require.NoError(t, err)

	buildOutput := `RUN apt-get update && \
 apt-get install -y curl && \
 curl -o helm.tgz https://get.helm.sh/helm-%s-linux-amd64.tar.gz && \
 tar -xzf helm.tgz && \
 mv linux-amd64/helm /usr/local/bin && \
 rm helm.tgz
RUN helm init --client-only
RUN apt-get update && \
 apt-get install -y apt-transport-https curl && \
 curl -o kubectl https://storage.googleapis.com/kubernetes-release/release/v1.15.3/bin/linux/amd64/kubectl && \
 mv kubectl /usr/local/bin && \
 chmod a+x /usr/local/bin/kubectl`

	t.Run("build with a valid config", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-valid-config.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")
		wantOutput := fmt.Sprintf(buildOutput, m.HelmClientVersion) +
			"\nRUN helm repo add stable kubernetes-charts" +
			"\nRUN helm repo update"
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build with a valid config and multiple repositories", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-valid-config-multi-repos.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")
		gotOutput := m.TestContext.GetOutput()
		assert.Contains(t, gotOutput, fmt.Sprintf(buildOutput, m.HelmClientVersion))
		assert.Contains(t, gotOutput, "RUN helm repo add harbor https://helm.getharbor.io")
		assert.Contains(t, gotOutput, "RUN helm repo add jetstack https://charts.jetstack.io")
		assert.Contains(t, gotOutput, "RUN helm repo add stable kubernetes-charts")
		assert.Contains(t, gotOutput, "RUN helm repo update")
	})

	t.Run("build with invalid config", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-invalid-config-url.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")
		wantOutput := fmt.Sprintf(buildOutput, m.HelmClientVersion) +
			"\nRUN helm repo update"
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build with a defined helm client version", func(t *testing.T) {

		b, err := ioutil.ReadFile("testdata/build-input-with-supported-client-version.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)
		err = m.Build()
		require.NoError(t, err, "build failed")
		wantOutput := fmt.Sprintf(buildOutput, m.HelmClientVersion)
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build with a defined helm client version that does not meet the semver constraint", func(t *testing.T) {

		b, err := ioutil.ReadFile("testdata/build-input-with-unsupported-client-version.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)
		err = m.Build()
		require.EqualError(t, err, `supplied clientVersion "v3.2.1" does not meet semver constraint "^v2.x"`)
	})

	t.Run("build with a defined helm client version that does not parse as valid semver", func(t *testing.T) {

		b, err := ioutil.ReadFile("testdata/build-input-with-invalid-client-version.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)
		err = m.Build()
		require.EqualError(t, err, `supplied client version "v3.2.1.0" cannot be parsed as semver: Invalid Semantic Version`)
	})
}
