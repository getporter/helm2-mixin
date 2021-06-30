package helm2

import (
	"fmt"
	"os"
	"testing"

	"get.porter.sh/porter/pkg/test"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMixin_Init_TillerNotReady(t *testing.T) {
	os.Setenv(test.ExpectedCommandEnv, "helm init --service-account=tiller-deploy --upgrade --wait")
	h := NewTestMixin(t)

	initer := NewMockTillerIniter()
	initer.GetTillerVersion = func(m *Mixin) (string, error) {
		return "", errors.New(tillerNotReadyErr)
	}
	h.Mixin.TillerIniter = initer

	err := h.Init()
	require.NoError(t, err)

	gotOutput := h.TestContext.GetOutput()
	wantOutput := "Tiller is not ready; attempting to init.\n"
	require.Equal(t, wantOutput, gotOutput)
}

func TestMixin_Init_FailedRBACSetup(t *testing.T) {
	os.Setenv(test.ExpectedCommandEnv, "helm init --service-account=tiller-deploy --upgrade --wait")
	h := NewTestMixin(t)

	initer := NewMockTillerIniter()
	initer.GetTillerVersion = func(m *Mixin) (string, error) {
		return "", errors.New(tillerNotReadyErr)
	}
	initer.SetupTillerRBAC = func(m *Mixin) error {
		return errors.New("failed to setup RBAC")
	}

	h.Mixin.TillerIniter = initer

	err := h.Init()
	require.EqualError(t, err, "failed to setup RBAC for Tiller: failed to setup RBAC")

	gotOutput := h.TestContext.GetOutput()
	wantOutput := "Tiller is not ready; attempting to init.\n"
	require.Equal(t, wantOutput, gotOutput)
}

func TestMixin_Init_MismatchedVersion(t *testing.T) {
	h := NewTestMixin(t)

	initer := NewMockTillerIniter()
	initer.GetTillerVersion = func(m *Mixin) (string, error) {
		return "mismatchedVersion", nil
	}
	h.Mixin.TillerIniter = initer

	err := h.Init()
	require.NoError(t, err)

	gotOutput := h.TestContext.GetOutput()
	wantOutput := fmt.Sprintf("Tiller version (mismatchedVersion) does not match client version (%s); downloading a compatible client.\n", h.HelmClientVersion)
	require.Equal(t, wantOutput, gotOutput)
}

func TestMixin_Init_FailedClientInstall(t *testing.T) {
	h := NewTestMixin(t)

	initer := NewMockTillerIniter()
	initer.GetTillerVersion = func(m *Mixin) (string, error) {
		return "mismatchedVersion", nil
	}
	initer.InstallHelmClient = func(m *Mixin, version string) error {
		return errors.New("failed to install helm client")
	}
	h.Mixin.TillerIniter = initer

	err := h.Init()
	require.EqualError(t, err, "unable to install a compatible helm client: failed to install helm client")

	gotOutput := h.TestContext.GetOutput()
	wantOutput := fmt.Sprintf("Tiller version (mismatchedVersion) does not match client version (%s); downloading a compatible client.\n", h.HelmClientVersion)
	require.Equal(t, wantOutput, gotOutput)
}
