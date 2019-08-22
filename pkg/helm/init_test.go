package helm

import (
	"os"
	"testing"

	"github.com/deislabs/porter/pkg/test"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type FailingTillerVersionGetter struct {
	MockTillerIniter
}

func (f FailingTillerVersionGetter) getTillerVersion(m *Mixin) (string, error) {
	return "", errors.New(tillerNotReadyErr)
}

func TestMixin_Init_TillerNotReady(t *testing.T) {
	os.Setenv(test.ExpectedCommandEnv, "helm init --service-account=tiller-deploy --upgrade --wait")
	h := NewTestMixin(t)

	initer := FailingTillerVersionGetter{MockTillerIniter: MockTillerIniter{}}
	h.Mixin.TillerIniter = initer

	err := h.Init()
	require.NoError(t, err)

	gotOutput := h.TestContext.GetOutput()
	wantOutput := "Tiller is not ready; attempting to init.\n"
	require.Equal(t, wantOutput, gotOutput)
}

type FailingRBACSetterUpper struct {
	MockTillerIniter
}

func (f FailingRBACSetterUpper) getTillerVersion(m *Mixin) (string, error) {
	return "", errors.New(tillerNotFoundErr)
}

func (f FailingRBACSetterUpper) setupTillerRBAC(m *Mixin) error {
	return errors.New("failed to setup RBAC")
}

func TestMixin_Init_FailedRBACSetup(t *testing.T) {
	os.Setenv(test.ExpectedCommandEnv, "helm init --service-account=tiller-deploy --upgrade --wait")
	h := NewTestMixin(t)

	initer := FailingRBACSetterUpper{MockTillerIniter: MockTillerIniter{}}
	h.Mixin.TillerIniter = initer

	err := h.Init()
	require.EqualError(t, err, "failed to setup RBAC for Tiller: failed to setup RBAC")

	gotOutput := h.TestContext.GetOutput()
	wantOutput := "Tiller is not ready; attempting to init.\n"
	require.Equal(t, wantOutput, gotOutput)
}

type MismatchedTillerVersionGetter struct {
	MockTillerIniter
}

func (f MismatchedTillerVersionGetter) getTillerVersion(m *Mixin) (string, error) {
	return "mismatchedVersion", nil
}

func TestMixin_Init_MismatchedVersion(t *testing.T) {
	h := NewTestMixin(t)

	initer := MismatchedTillerVersionGetter{MockTillerIniter: MockTillerIniter{}}
	h.Mixin.TillerIniter = initer

	err := h.Init()
	require.NoError(t, err)

	gotOutput := h.TestContext.GetOutput()
	wantOutput := "Tiller version (mismatchedVersion) does not match client version (v2.14.3); downloading a compatible client.\n"
	require.Equal(t, wantOutput, gotOutput)
}

type FailingHelmClientFetcher struct {
	MockTillerIniter
}

func (f FailingHelmClientFetcher) getTillerVersion(m *Mixin) (string, error) {
	return "mismatchedVersion", nil
}

func (f FailingHelmClientFetcher) installHelmClient(m *Mixin, version string) error {
	return errors.New("failed to install helm client")
}

func TestMixin_Init_FailedClientInstall(t *testing.T) {
	h := NewTestMixin(t)

	initer := FailingHelmClientFetcher{MockTillerIniter: MockTillerIniter{}}
	h.Mixin.TillerIniter = initer

	err := h.Init()
	require.EqualError(t, err, "unable to install a compatible helm client: failed to install helm client")

	gotOutput := h.TestContext.GetOutput()
	wantOutput := "Tiller version (mismatchedVersion) does not match client version (v2.14.3); downloading a compatible client.\n"
	require.Equal(t, wantOutput, gotOutput)
}
