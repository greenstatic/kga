package app

import (
	"fmt"
	"github.com/greenstatic/kga/pkg/files"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

const (
	ConfigKind    = "kga-app"
	ConfigVersion = "v1alpha"
)

type Config struct {
	Kind    string `yaml:"kind"`
	Version string `yaml:"version"`
	Name    string `yaml:"name"`
	Spec    *Spec  `yaml:"spec,omitempty"`
}

type Spec struct {
	Namespace string             `yaml:"namespace,omitempty"`
	Type      string             `yaml:"type"`
	Helm      *HelmSpec          `yaml:"helm,omitempty"`
	Manifest  *ManifestSpec      `yaml:"manifest,omitempty"`
	Exclude   *[]ExcludeItemSpec `yaml:"exclude,omitempty"`
}

type HelmSpec struct {
	ChartName  string `yaml:"chartName"`
	Version    string `yaml:"version"`
	RepoName   string `yaml:"repoName"`
	RepoUrl    string `yaml:"repoUrl"`
	ValuesFile string `yaml:"valuesFile"`
}

type ManifestSpec struct {
	Version  string            `yaml:"version"`
	Urls     []string          `yaml:"urls"`
	Template map[string]string `yaml:"template,omitempty"`
}

type ExcludeItemSpec map[interface{}]interface{}

func ConfigYaml(c *Config) (string, error) {
	s, err := yaml.Marshal(*c)
	if err != nil {
		return string(s), errors.Wrap(err, "failed to marshal YAML")
	}
	return string(s), nil
}

func PathIsKgaApp(path string) (bool, error) {
	// TODO - check if kga.yaml appears to be a valid kga app configuration
	return files.FileOrDirExists(filepath.Join(path, "kga.yaml"))
}

type ConfigBadFieldValueError string
func (s ConfigBadFieldValueError) Error() string {
	return fmt.Sprintf("kga config bad field value: %s", string(s))
}

type ConfigMissingFieldError string
func (s ConfigMissingFieldError) Error() string {
	return fmt.Sprintf("kga config missing field: %s", string(s))
}

func (c *Config) Verify() error {
	if c.Kind != ConfigKind {
		return ConfigBadFieldValueError(".kind")
	}

	if c.Name == "" {
		return ConfigBadFieldValueError(".name cannot be empty")
	}

	if c.Version != ConfigVersion {
		return ConfigBadFieldValueError(".version is unknown")
	}

	if c.Spec == nil {
		return ConfigMissingFieldError(".spec")
	}

	if c.Spec.Type == "" {
		return ConfigMissingFieldError(".spec.type")
	}

	switch c.Spec.Type {
	case BasicType:
	case ManifestType:
		if c.Spec.Manifest == nil {
			return ConfigMissingFieldError(".spec.manifest")
		}

		if c.Spec.Helm != nil {
			return ConfigBadFieldValueError(".spec.helm is defined for type: manifest")
		}

		if c.Spec.Manifest.Version == "" {
			return ConfigMissingFieldError(".spec.manifest.version")
		}

		if len(c.Spec.Manifest.Urls) == 0 {
			return ConfigMissingFieldError(".spec.manifest.urls")
		}

	case HelmType:
		if c.Spec.Helm == nil {
			return ConfigMissingFieldError(".spec.helm")
		}

		if c.Spec.Manifest != nil {
			return ConfigBadFieldValueError(".spec.manifest is defined for type: helm")
		}

		if c.Spec.Namespace == "" {
			return ConfigMissingFieldError(".spec.namespace")
		}

		if c.Spec.Helm.ChartName == "" {
			return ConfigMissingFieldError(".spec.helm.chartName")
		}
		if c.Spec.Helm.Version == "" {
			return ConfigMissingFieldError(".spec.helm.version")
		}
		if c.Spec.Helm.RepoName == "" {
			return ConfigMissingFieldError(".spec.helm.repoName")
		}
		if c.Spec.Helm.RepoUrl == "" {
			return ConfigMissingFieldError(".spec.helm.repoUrl")
		}

	default:
		return ConfigBadFieldValueError(".spec.type can be only: basic, manifest or helm")
	}

	return nil
}

func ParseFromFileConfig(filepath string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read contents of: %s", filepath)
	}

	c := Config{}
	return &c, yaml.Unmarshal(data, &c)
}