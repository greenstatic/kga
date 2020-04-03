package layout

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func CreateBaseManifestsDir(appPath string) error {
	path := filepath.Join(appPath, "base", "manifests")
	return os.MkdirAll(path, os.FileMode(0755))
}

func CreateBaseKustomization(appPath string) error {
	manifestResources, err := enumerateAllBaseManifestFiles(appPath)
	if err != nil {
		return err
	}

	kustomizationData, err := baseKustomization(manifestResources)
	if err != nil {
		return err
	}

	return baseKustomizationWrite(appPath, kustomizationData)
}

func baseKustomizationWrite(appPath string, data []byte) error {
	path := filepath.Join(appPath, "base", "kustomization.yaml")
	return ioutil.WriteFile(path, data, os.FileMode(0640))
}

func baseKustomization(manifestResources []string) ([]byte, error) {

	templateStr := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
{{ range $val := . -}}
- manifests/{{ $val }}
{{- end }}
`

	tmpl := template.Must(template.New("kustomization").Parse(templateStr))

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, manifestResources); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func enumerateAllBaseManifestFiles(appName string) ([]string, error) {
	path := filepath.Join(appName, "base", "manifests")

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	files2 := make([]string, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			files2 = append(files2, f.Name())
		}
	}

	return files2, nil
}

func CreateMainKustomizationFile(appPath string) error {
	exists, err := FileOrDirExists(filepath.Join(appPath, "kustomization.yaml"))
	if err != nil {
		return err
	}
	if !exists {
		return createMainKustomizationFile(appPath)
	}
	return nil
}

func createMainKustomizationFile(appPath string) error {
	path := filepath.Join(appPath, "kustomization.yaml")

	contents := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
  - overlay
`

	return ioutil.WriteFile(path, []byte(contents), os.FileMode(0640))
}
