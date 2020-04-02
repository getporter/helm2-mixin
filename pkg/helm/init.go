package helm

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

const tillerNotReadyErr string = "could not find a ready tiller pod"

const tillerNotFoundErr string = "could not find tiller"

// TillerIniter is an interface for methods associated with Tiller interactions
type TillerIniter interface {
	setupTillerRBAC(m *Mixin) error
	getTillerVersion(m *Mixin) (string, error)
	runRBACResourceCmd(m *Mixin, cmd *exec.Cmd) error
	installHelmClient(m *Mixin, version string) error
}

// RealTillerIniter implements the TillerIniter interface, in REAL life
type RealTillerIniter struct{}

// Init inits the Helm server (Tiller) if not running, else, ensures the
// Helm client matches the installed Tiller version
func (m *Mixin) Init() error {
	ti := m.TillerIniter

	tillerVersion, err := ti.getTillerVersion(m)
	if err != nil {
		switch errMsg := err.Error(); {
		case strings.Contains(errMsg, tillerNotReadyErr) || strings.Contains(errMsg, tillerNotFoundErr):
			fmt.Fprintln(m.Out, "Tiller is not ready; attempting to init.")

			err := ti.setupTillerRBAC(m)
			if err != nil {
				return errors.Wrap(err, "failed to setup RBAC for Tiller")
			}

			initCmd := m.NewCommand("helm", "init", "--service-account=tiller-deploy", "--upgrade", "--wait")
			prettyCmd := fmt.Sprintf("%s %s", initCmd.Path, strings.Join(initCmd.Args, " "))

			initCmd.Stdout = m.Out
			initCmd.Stderr = m.Err

			err = initCmd.Start()
			if err != nil {
				return errors.Wrapf(err, "could not execute command, %s", prettyCmd)
			}
			err = initCmd.Wait()
			if err != nil {
				return errors.Wrap(err, "unable to init Tiller")
			}
		default:
			return errors.Wrap(err, "unable to communicate with Tiller")
		}
	} else {
		if m.getHelmClientVersion() != tillerVersion {
			fmt.Fprintf(m.Out, "Tiller version (%s) does not match client version (%s); downloading a compatible client.\n",
				tillerVersion, m.getHelmClientVersion())

			err := ti.installHelmClient(m, tillerVersion)
			if err != nil {
				return errors.Wrap(err, "unable to install a compatible helm client")
			}
		}
	}
	return nil
}

func (r RealTillerIniter) setupTillerRBAC(m *Mixin) error {
	cmd := m.NewCommand("kubectl", "create", "serviceaccount", "-n", "kube-system", "tiller-deploy")
	err := r.runRBACResourceCmd(m, cmd)
	if err != nil {
		return err
	}

	cmd = m.NewCommand("kubectl", "create", "clusterrolebinding", "tiller-deploy",
		"--clusterrole", "cluster-admin", "--serviceaccount", "kube-system:tiller-deploy")
	return r.runRBACResourceCmd(m, cmd)
}

func (r RealTillerIniter) runRBACResourceCmd(m *Mixin, cmd *exec.Cmd) error {
	var stderr bytes.Buffer

	prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
	cmd.Stdout = m.Out
	// We'll be checking stderr to determine whether or not to error out or ignore
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "could not execute command, %s", prettyCmd)
	}
	err = cmd.Wait()
	if err != nil {
		// Only return the error if other than a pre-existence error
		if !strings.Contains(stderr.String(), "already exists") {
			return errors.Wrapf(err,
				"unable to create RBAC resource: %s", stderr.String())
		}
	}
	return nil
}

func (r RealTillerIniter) getTillerVersion(m *Mixin) (string, error) {
	var stderr bytes.Buffer

	cmd := m.NewCommand("helm", "version", "--server")
	cmd.Stderr = &stderr

	outputBytes, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "unable to determine Helm's server version: %s", stderr.String())
	}
	re := regexp.MustCompile(`v[0-9]*\.[0-9]*\.[0-9]*`)
	version := re.FindString(string(outputBytes))

	return version, nil
}

func (r RealTillerIniter) installHelmClient(m *Mixin, version string) error {
	helmArchive := fmt.Sprintf(helmArchiveTmpl, version)
	url := fmt.Sprintf(helmDownloadURLTmpl, helmArchive)

	// Fetch archive from url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to construct GET request for fetching helm client binary")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to download helm client binary via url: %s", url)
	}
	defer res.Body.Close()

	// Create a temp dir
	tmpDir, err := m.FileSystem.TempDir("", "tmp")
	if err != nil {
		return errors.Wrap(err, "unable to create a temporary directory for downloading the helm client binary")
	}
	defer os.RemoveAll(tmpDir)

	// Create the local archive
	archiveFile, err := m.FileSystem.Create(filepath.Join(tmpDir, helmArchive))
	if err != nil {
		return errors.Wrap(err, "unable to create a local file for the helm client binary")
	}

	// Copy response body to local archive
	_, err = io.Copy(archiveFile, res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to copy the helm client binary to the local archive file")
	}

	// Create a dir for unarchived contents
	unarchivedDir, err := m.FileSystem.TempDir(tmpDir, "helm")
	if err != nil {
		return errors.Wrap(err, "unable to create a temporary directory for the unarchived helm client download contents")
	}

	err = archiver.Unarchive(archiveFile.Name(), unarchivedDir)
	if err != nil {
		return errors.Wrap(err, "unable to unarchive the helm client download")
	}

	// Move the helm binary into the appropriate location
	binPath := "/usr/local/bin/helm"
	err = m.FileSystem.Rename(fmt.Sprintf("%s/linux-amd64/helm", unarchivedDir), binPath)
	if err != nil {
		return errors.Wrapf(err, "unable to install the helm client binary to %q", binPath)
	}
	return nil
}
