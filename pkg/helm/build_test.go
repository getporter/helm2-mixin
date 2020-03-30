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
	const HelmDefaultClientVersion string = "v2.15.2"
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
		wantOutput := fmt.Sprintf(buildOutput, HelmDefaultClientVersion) + "\nRUN helm repo add stable kubernetes-charts --username username --password password"
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build with invalid config", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-invalid-config.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")
		wantOutput := fmt.Sprintf(buildOutput, HelmDefaultClientVersion)
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build with a defined helm client version", func(t *testing.T) {
		var version = "v2.16.1"
		b, err := ioutil.ReadFile("testdata/build-input-with-version.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)
		err = m.Build()
		require.NoError(t, err, "build failed")
		wantOutput := fmt.Sprintf(buildOutput, version)
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
		helmClientVersion = HelmDefaultClientVersion
	})
}
