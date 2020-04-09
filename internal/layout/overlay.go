package layout

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func OverlayExists(appPath string) (bool, error) {
	path := filepath.Join(appPath, "overlay")
	return FileOrDirExists(path)
}

func OverlayCreateGeneralLayout(appPath, namespace string, overrideNamespaceInKustomization bool) error {
	if err := overlayCreateDirectoryLayout(appPath); err != nil {
		return err
	}

	tmpNamespace := namespace
	if !overrideNamespaceInKustomization {
		tmpNamespace = ""
	}
	if err := overlayCreateKustomizationFile(appPath, tmpNamespace); err != nil {
		return err
	}
	if err := overlayCreateNamespaceResourceFile(appPath, namespace); err != nil {
		return err
	}

	return nil
}

func overlayCreateDirectoryLayout(appPath string) error {
	dirs := []string{
		"overlay/patches",
		"overlay/resources",
	}

	for _, val := range dirs {
		if err := os.MkdirAll(filepath.Join(appPath, val), os.FileMode(0775)); err != nil {
			return err
		}
	}

	return nil
}

func overlayCreateKustomizationFile(appPath, namespace string) error {
	path := filepath.Join(appPath, "overlay", "kustomization.yaml")

	templateStr := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
{{ if . }}namespace: {{ . }}{{ end }}
bases:
  - ../base

resources:
- resources/namespace.yaml

#patchesStrategicMerge:
#  - patches/patch.yaml
`

	templ := template.Must(template.New("namespace").Parse(templateStr))
	buff := new(bytes.Buffer)
	if err := templ.Execute(buff, namespace); err != nil {
		return err
	}


	return ioutil.WriteFile(path, buff.Bytes(), os.FileMode(0640))
}

func overlayCreateNamespaceResourceFile(appPath, namespace string) error {
	path := filepath.Join(appPath, "overlay", "resources", "namespace.yaml")

	templateStr := `apiVersion: v1
kind: Namespace
metadata:
  name: {{ . }}
`

	templ := template.Must(template.New("namespace").Parse(templateStr))
	buff := new(bytes.Buffer)
	if err := templ.Execute(buff, namespace); err != nil {
		return err
	}

	return ioutil.WriteFile(path, buff.Bytes(), os.FileMode(0640))
}
