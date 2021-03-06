package app

import (
	"github.com/greenstatic/kga/pkg/files"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Generic initialization function for any app type
func Init(t Type, path, name string) error {
	// Create base app dir
	if err := os.MkdirAll(path, os.FileMode(0775)); err != nil {
		return err
	}

	c := Config{Kind: ConfigKind, Version: ConfigVersion, Name: name}
	t.initConfig(&c)

	if c.Spec.Exclude == nil {
		c.Spec.Exclude = new([]ExcludeItemSpec)
	}
	// Add default "kind: Secret" exclude
	excludeSecret := ExcludeItemSpec{}
	excludeSecret["kind"] = "Secret"
	*c.Spec.Exclude = append(*c.Spec.Exclude, excludeSecret)

	if err := saveKgaYaml(&c, filepath.Join(path, "kga.yaml")); err != nil {
		return errors.Wrap(err, "failed to save kga.yaml file")
	}

	if err := createInitAppStructureBaseDir(path); err != nil {
		return err
	}
	if err := createInitAppStructureOverlay(path); err != nil {
		return err
	}
	if err := createInitAppStructureRootKustomization(path); err != nil {
		return err
	}

	// app type specific initialization
	if err := t.init(&c, path); err != nil {
		return err
	}

	return nil
}

func saveKgaYaml(c *Config, filepath string) error {
	s, err := ConfigYaml(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, []byte(s), os.FileMode(0644))
}

func createInitAppStructureBaseDir(path string) error {
	return os.MkdirAll(filepath.Join(path, "base"), os.FileMode(0775))
}

func createInitAppStructureBase(path string) error {
	dirs := []string{"base/manifests"}
	if err := files.CreteDirs(dirs, path); err != nil {
		return errors.Wrap(err, "failed to create app directory structure")
	}

	baseKustomization := `# This file has been generated by kga
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
#resources:
#- manifests/example.yaml
`
	if err := ioutil.WriteFile(filepath.Join(path, "base/kustomization.yaml"), []byte(baseKustomization),
		os.FileMode(0644)); err != nil {
		return errors.Wrap(err, "filed to createfile: base/kustomization.yaml")
	}

	return nil
}

func createInitAppStructureOverlay(path string) error {
	dirs := []string{"overlay/patches", "overlay/resources"}
	if err := files.CreteDirs(dirs, path); err != nil {
		return errors.Wrap(err, "failed to create app directory structure")
	}

	overlayKustomization := `# This file has been generated by kga
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../base
#  - resources/example.yaml

#patchesJson6902:
#  - target:
#      name: example-ingress
#      group: networking.k8s.io
#      version: v1beta1
#      kind: Ingress
#    path: patches/example_json.yaml

#patchesStrategicMerge:
#  - patches/example.yaml
`
	if err := ioutil.WriteFile(filepath.Join(path, "overlay/kustomization.yaml"), []byte(overlayKustomization),
		os.FileMode(0644)); err != nil {
		return errors.Wrap(err, "filed to create file: overlay/kustomization.yaml")
	}

	return nil
}

func createInitAppStructureRootKustomization(path string) error {
	mainKustomization := `# This file has been generated by kga
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - overlay
`

	if err := ioutil.WriteFile(filepath.Join(path, "kustomization.yaml"), []byte(mainKustomization),
		os.FileMode(0644)); err != nil {
		return errors.Wrap(err, "filed to create file: kustomization.yaml")
	}

	return nil
}
