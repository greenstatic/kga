package app

import (
	"bytes"
	"fmt"
	"github.com/greenstatic/kga/pkg/files"
	"github.com/greenstatic/kga/pkg/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

const ManifestType = "manifest"

type Manifest struct {
}

func (_ Manifest) AppType() string {
	return ManifestType
}

func (m *Manifest) init(c *Config, path string) error {
	return nil
}

func (m *Manifest) initConfig(c *Config) {
	c.Spec = &Spec{Type: ManifestType, Namespace: c.Name}

	c.Spec.Manifest = &ManifestSpec{
		Version:  "# TODO",
		Urls:     []string{"# TODO"},
		Template: nil,
	}
}

func (m *Manifest) generate(c *Config, path string) error {
	// Delete base dir
	if err := files.RemoveDirIfExists(filepath.Join(path, "base")); err != nil {
		return err
	}

	// Create new base dir
	if err := m.createGenerateAppStructureBase(path); err != nil {
		return err
	}

	if err := downloadAndSaveManifests(c, filepath.Join(path, "base", "manifests")); err != nil {
		return err
	}

	if c.Spec.Exclude != nil {
		if err := removeExcludedResources(filepath.Join(path, "base", "manifests"),
			filepath.Join(path, "base", "excluded"), c.Spec.Exclude); err != nil {
			return err
		}
	}

	if err := generateBase(filepath.Join(path, "base")); err != nil {
		return err
	}

	// Overlay
	if c.Spec.Namespace != "" {
		// Namespace is defined, remove namespace resource from base manifest and
		// redefine it in overlay along with a namespace override in overlay/kustomization.yaml

		// TODO


	}

	return nil
}

func downloadAndSaveManifests(c *Config, manifestPath string) error {
	for _, urlTemplate := range c.Spec.Manifest.Urls {
		url, err := manifestUrlApplyTemplate(c, urlTemplate)
		if err != nil {
			return errors.Wrap(err, "URL templating failed")
		}

		log.Info("Fetching content from URL: " + url)
		contents, err := downloadUrlContents(url)
		if err != nil {
			return errors.Wrap(err, "failed to download URL contents")
		}

		filename := urlFileName(url)
		path := filepath.Join(manifestPath, filename)

		log.Infof("Saving file %s", path)
		if isError := saveContents(path, contents); isError != nil {
			return errors.Wrap(err, "failed to save downloaded manifest contents to file")
		}
	}

	return nil
}

type manifestTemplateFields struct {
	Version string
	Template map[string]string
	Config *Config
}

func manifestUrlApplyTemplate(c *Config, url string) (string, error) {
	t, err := template.New("url").Parse(url)
	if err != nil {
		return "", err
	}

	fields := manifestTemplateFields{
		Version:  c.Spec.Manifest.Version,
		Template: c.Spec.Manifest.Template,
		Config:   c,
	}

	buff := new(bytes.Buffer)
	if err := t.Execute(buff, fields); err != nil {
		return "", err
	}

	return buff.String(), nil
}

func downloadUrlContents(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode / 100 != 2 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "failed to read %d HTTP status code body contents", resp.StatusCode)
		}
		log.Error(string(b))
		return []byte{}, errors.New(fmt.Sprintf("HTTP status code: %d", resp.StatusCode))
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to read contents of body")
	}

	return buff, nil
}

func urlFileName(url string) string {
	return filepath.Base(url)
}

func saveContents(filepath string, contents []byte) error {
	return ioutil.WriteFile(filepath, contents, os.FileMode(0640))
}

func (m *Manifest) createGenerateAppStructureBase(path string) error {
	dirs := []string{"base/manifests"}
	if err := files.CreteDirs(dirs, path); err != nil {
		return errors.Wrap(err, "failed to create app directory structure")
	}
	return nil
}