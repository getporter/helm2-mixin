package helm

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

// Init inits the Helm server (Tiller) if not running, else, ensures the
// Helm client matches the installed Tiller version
func (m *Mixin) Init() error {
	tillerVersion, err := m.getTillerVersion()
	if err != nil {
		fmt.Fprintf(m.Out, "Tiller may not be running; attempting to init.")
		initCmd := m.NewCommand("helm", "init", "--upgrade", "--wait")
		prettyCmd := fmt.Sprintf("%s %s", initCmd.Path, strings.Join(initCmd.Args, " "))

		initCmd.Stdout = m.Out
		initCmd.Stderr = m.Err

		err = initCmd.Start()
		if err != nil {
			return fmt.Errorf("could not execute command, %s: %s", prettyCmd, err)
		}
		err = initCmd.Wait()
		if err != nil {
			return errors.Wrap(err, "unable to init Tiller")
		}
		// TODO: ughhhh RBAC setup may be needed?
	} else {
		if helmClientVersion != tillerVersion {
			fmt.Fprintf(m.Out, "Tiller version (%s) does not match client version (%s); downloading a compatible client.\n",
				tillerVersion, helmClientVersion)

			err := installHelmClient(tillerVersion)
			if err != nil {
				return errors.Wrap(err, "unable to install a compatible helm client")
			}
		}
	}
	return nil
}

func (m *Mixin) getTillerVersion() (string, error) {
	cmd := m.NewCommand("helm", "version", "--server")
	outputBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`v[0-9]*\.[0-9]*\.[0-9]*`)
	version := re.FindString(string(outputBytes))

	return version, nil
}

func installHelmClient(version string) error {
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
	tmpDir, err := ioutil.TempDir("", "tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Create the local archive
	archiveFile, err := os.Create(filepath.Join(tmpDir, helmArchive))
	if err != nil {
		return err
	}

	// Copy response body to local archive
	_, err = io.Copy(archiveFile, res.Body)
	if err != nil {
		return err
	}

	// Create a dir for unarchived contents
	unarchivedDir, err := ioutil.TempDir(tmpDir, "helm")
	if err != nil {
		return err
	}

	err = archiver.Unarchive(archiveFile.Name(), unarchivedDir)
	if err != nil {
		return err
	}

	// Move the helm binary into the appropriate location
	err = os.Rename(fmt.Sprintf("%s/linux-amd64/helm", unarchivedDir), "/usr/local/bin/helm")
	if err != nil {
		return err
	}
	return nil
}
