package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

func ParseFile(path string) (*Config, error) {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read file")
	}

	c := Config{}

	err = yaml.Unmarshal(fileContents, &c)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse config")
	}

	return &c, nil
}

func VerifyKgaFile(path string) error {
	conf, err := ParseFile(path)
	if err != nil {
		return err
	}

	if conf.Spec != nil {
		if conf.Spec.Namespace == "" {
			return errors.New("spec.namespace field is missing")
		}

		if conf.Spec.Helm == nil && conf.Spec.Manifest == nil {
			return errors.New("spec.helm or spec.manifest field required (exclusive disjunction)")
		}

		// Check for mutually exclusive helm and manifest fields
		if !existsNil(conf.Spec.Helm, conf.Spec.Manifest) {
			return errors.New("spec.helm and spec.manifest fields are mutually exclusive - only one is allowed")
		}

		// Check all helm fields are present
		if conf.Spec.Helm != nil {
			if conf.Spec.Helm.Version == "" {
				return errors.New("spec.helm.version field is missing")
			}
			if conf.Spec.Helm.RepoName == "" {
				return errors.New("spec.helm.repoName field is missing")
			}
			if conf.Spec.Helm.RepoUrl == "" {
				return errors.New("spec.helm.repoUrl field is missing")
			}
			if conf.Spec.Helm.ChartName == "" {
				return errors.New("spec.helm.chartName field is missing")
			}
		}

		if conf.Spec.Manifest != nil {
			if len(conf.Spec.Manifest.Urls) == 0 {
				return errors.New("spec.manifest.urls is empty")
			}
		}

		// Check all Manifest contains at least one element if it is specified
		if conf.Spec.Manifest != nil && len(conf.Spec.Manifest.Urls) == 0 {
			return errors.New("spec.manifest.Urls does not contain any urls")
		}
	} else {
		return errors.New("spec field missing")
	}

	return nil
}

func existsNil(x ...interface{}) bool {
	for _, val := range x {
		if val == nil || (reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()) {
			return true
		}
	}
	return false
}
