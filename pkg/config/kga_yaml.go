package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func KgaYamlCreate(appType AppType, appName, appPath string) error {
	buf, err := kgaYamlContents(appType, appName)
	if err != nil {
		return err
	}
	return kgaYamlFileCreate(appPath, buf)
}

func kgaYamlFileCreate(appPath string, content *bytes.Buffer) error {
	path := filepath.Join(appPath, "kga.yaml")
	return ioutil.WriteFile(path, content.Bytes(), os.FileMode(0644))
}

func kgaYamlContents(appType AppType, appName string) (*bytes.Buffer, error) {
	templateStr := ""
	switch appType {
	case AppTypeHelm:
		templateStr = kgaHelmTemplate
	case AppTypeManifest:
		templateStr = kgaManifestTemplate
	default:
		return nil, errors.New("unknown AppType")
	}

	templateVar := map[string]string{
		"appName":    appName,
		"kgaVersion": configVersion,
	}

	tmpl := template.Must(template.New("config").Parse(templateStr))
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, templateVar); err != nil {
		return nil, err
	}

	return buf, nil
}

var kgaHelmTemplate = `kind: kga-app
version: {{ .kgaVersion }}
name: {{ .appName }}
spec:
  helm:
    chartName: # TODO
    version:   # TODO
    repoName:  # TODO
    repoUrl:   # TODO
    namespace: {{ .appName }}
    valuesFile: ./helm_values.yaml
`

var kgaManifestTemplate = `kind: kga-app
version: {{ .kgaVersion }}
name: {{ .appName }}
spec:
  manifest:
    version: v2.0.0
    urls:
    - "https://example.com/{{"{{"}} .Version {{"}}"}}/manifests.yaml" # TODO - replace
`
