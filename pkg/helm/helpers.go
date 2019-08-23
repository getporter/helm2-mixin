package helm

import (
	"os/exec"
	"testing"

	"github.com/deislabs/porter/pkg/context"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
)

type TestMixin struct {
	*Mixin
	TestContext *context.TestContext
}

type testKubernetesFactory struct {
}

func (t *testKubernetesFactory) GetClient(configPath string) (kubernetes.Interface, error) {
	return testclient.NewSimpleClientset(), nil
}

type MockTillerIniter struct {
	GetTillerVersion   func(m *Mixin) (string, error)
	SetupTillerRBAC    func(m *Mixin) error
	RunRBACResourceCmd func(m *Mixin, cmd *exec.Cmd) error
	InstallHelmClient  func(m *Mixin, version string) error
}

func (t MockTillerIniter) getTillerVersion(m *Mixin) (string, error) {
	return t.GetTillerVersion(m)
}

func (t MockTillerIniter) setupTillerRBAC(m *Mixin) error {
	return t.SetupTillerRBAC(m)
}

func (t MockTillerIniter) runRBACResourceCmd(m *Mixin, cmd *exec.Cmd) error {
	return t.RunRBACResourceCmd(m, cmd)
}

func (t MockTillerIniter) installHelmClient(m *Mixin, version string) error {
	return t.InstallHelmClient(m, version)
}

func NewMockTillerIniter() MockTillerIniter {
	return MockTillerIniter{
		GetTillerVersion: func(m *Mixin) (string, error) {
			return helmClientVersion, nil
		},
		SetupTillerRBAC: func(m *Mixin) error {
			return nil
		},
		RunRBACResourceCmd: func(m *Mixin, cmd *exec.Cmd) error {
			return nil
		},
		InstallHelmClient: func(m *Mixin, version string) error {
			return nil
		},
	}
}

// NewTestMixin initializes a helm mixin, with the output buffered, and an in-memory file system.
func NewTestMixin(t *testing.T) *TestMixin {
	c := context.NewTestContext(t)
	m := New()
	m.Context = c.Context
	m.ClientFactory = &testKubernetesFactory{}
	m.TillerIniter = NewMockTillerIniter()
	return &TestMixin{
		Mixin:       m,
		TestContext: c,
	}
}
