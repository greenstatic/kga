package generate

import (
	"bytes"
	"github.com/greenstatic/kga/internal/layout"
	"github.com/greenstatic/kga/pkg/config"
	"github.com/greenstatic/kga/pkg/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func DownloadManifestFiles(appPath string, spec *config.ManifestSpec, cnfg *config.Config, removeNamespaceResources bool) {
	if err := layout.CreateBaseManifestsDir(appPath); err != nil {
		log.Error(err)
		log.Fatal("Failed to create base manifest directory")
	}

	for _, urlTemplate := range spec.Urls {
		url, err := manifestUrlApplyTemplate(urlTemplate, spec.Version, cnfg, spec.Template)
		if err != nil {
			log.Error(err)
			log.Fatal("URL templating failed")
		}

		log.Info("Fetching content from URL: " + url)
		contents, err := downloadUrlContents(url)
		if err != nil {
			log.Error(err)
			log.Fatal("Failed to download URL contents")
		}

		removedManifest := make([]byte, 0)

		if removeNamespaceResources {
			log.Info("Excluding namespace resources from manifest")
			removedManifestsStr, manifests, err := removeNamespaceResource(string(contents))
			if err != nil {
				log.Error(err)
				log.Fatal("Failed to remove namespace resources")
			}

			contents = []byte(manifests)
			removedManifest = []byte(removedManifestsStr)
		}

		filename := urlFileName(url)
		log.Infof("Saving manifests to base/manifests/%s", filename)
		if err := saveContentsToBaseManifests(appPath, filename, contents); err != nil {
			log.Error(err)
			log.Fatal("Failed to save downloaded manifest contents to file")
		}

		if len(removedManifest) > 0 {
			log.Infof("Saving excluded resources to base/excluded/%s", filename)
			if err := saveContentsToExcludedManifests(appPath, filename, removedManifest); err != nil {
				log.Error(err)
				log.Fatal("Failed to save excluded manifest contents to file")
			}
		}
	}
}

type templateFields struct {
	Version string
	Template map[string]string
	Config *config.Config
}

func manifestUrlApplyTemplate(url string, version string, cnfg *config.Config, tmpl map[string]string) (string, error) {
	t, err := template.New("url").Parse(url)
	if err != nil {
		return "", err
	}

	fields := templateFields{
		Version:  version,
		Template: tmpl,
		Config:   cnfg,
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
			return []byte{}, errors.Wrap(err, "failed to read 2xx HTTP status code body contents")
		}
		log.Error(string(b))
		return []byte{}, errors.New("non 2xx HTTP status code")
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

func saveContentsToBaseManifests(appPath string, filename string, contents []byte) error {
	path := filepath.Join(appPath, "base", "manifests")
	filePath := filepath.Join(path, filename)
	return ioutil.WriteFile(filePath, contents, os.FileMode(0640))
}

func saveContentsToExcludedManifests(appPath string, filename string, contents []byte) error {
	path := filepath.Join(appPath, "base", "excluded")
	filePath := filepath.Join(path, filename)

	if exists, err := layout.FileOrDirExists(path); err != nil || !exists {
		if exists {
			log.Info("Directory base/excluded exists")
		} else if err != nil {
			log.Info("Error while creating excluded manifest")
			return err
		} else {
			log.Infof("Creating %s", path)
			return os.Mkdir(path, os.FileMode(0775))
		}
	}

	return ioutil.WriteFile(filePath, contents, os.FileMode(0640))
}

func removeNamespaceResource(manifest string) (string, string, error) {
	excludeStr := `apiVersion: v1
kind: Namespace
`
	return excludeResourceFromManifest(excludeStr, manifest)
}