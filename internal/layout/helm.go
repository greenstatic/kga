package layout

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateHelmValuesFile(appPath string) error {
	path := filepath.Join(appPath, "helm_values.yaml")
	return ioutil.WriteFile(path, []byte{}, os.FileMode(0644))
}

func CreateBaseHelmManifests(appPath, chartName string, manifests []byte) error {
	if err := CreateBaseManifestsDir(appPath); err != nil {
		return errors.Wrap(err, "failed to create base/manifests dir")
	}

	if err := writeBaseHelmManifests(appPath, chartName, manifests); err != nil {
		return errors.Wrap(err, "failed to write manifest file")
	}

	return nil
}

func writeBaseHelmManifests(appPath, chartName string, manifests []byte) error {
	path := filepath.Join(appPath, "base", "manifests", chartName+".yaml")
	return ioutil.WriteFile(path, manifests, os.FileMode(0640))
}
