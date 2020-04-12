package app

import (
	"bytes"
	"fmt"
	"github.com/greenstatic/kga/pkg/files"
	"github.com/greenstatic/kga/pkg/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const HelmType = "helm"

type Helm struct {
}

func (_ Helm) AppType() string {
	return HelmType
}

func (h *Helm) init(c *Config, path string) error {
	helmValuesFile := filepath.Join(path, c.Spec.Helm.ValuesFile)
	if err := createHelmValuesFile(helmValuesFile, ""); err != nil {
		return errors.Wrapf(err, "failed to write helm values file: %s", helmValuesFile)
	}

	return nil
}

func (h *Helm) initConfig(c *Config) {
	c.Spec = &Spec{Type: HelmType, Namespace: c.Name}

	c.Spec.Helm = &HelmSpec{
		ChartName:  "# TODO",
		Version:    "# TODO",
		RepoName:   "# TODO",
		RepoUrl:    "# TODO",
		ValuesFile: "helm_values.yaml",
	}

}

func createHelmValuesFile(filepath string, content string) error {
	return ioutil.WriteFile(filepath, []byte(content), os.FileMode(0644))
}


func (h *Helm) generate(c *Config, path string) error {
	// Delete base dir
	if err := files.RemoveDirIfExists(filepath.Join(path, "base")); err != nil {
		return err
	}

	// Create new base dir
	if err := h.createGenerateAppStructureBase(path); err != nil {
		return err
	}

	if err := createHelmChartManifests(c, path); err != nil {
		return err
	}

	if err := generateBase(filepath.Join(path, "base")); err != nil {
		return err
	}

	if c.Spec.Exclude != nil {
		if err := removeExcludedResources(filepath.Join(path, "base", "manifests"),
			filepath.Join(path, "base", "excluded"), c.Spec.Exclude); err != nil {
			return err
		}
	}

	return nil
}

func (h *Helm) createGenerateAppStructureBase(path string) error {
	dirs := []string{"base/manifests"}
	if err := files.CreteDirs(dirs, path); err != nil {
		return errors.Wrap(err, "failed to create app directory structure")
	}
	return nil
}


func helmExecutablePath() string {
	envHelm := os.Getenv("HELM")
	if envHelm != "" {
		return envHelm
	}

	return "helm" // assume that it is in the users path
}

func helmCheckExecutable() error {
	cmdPath := helmExecutablePath()
	helm := exec.Command(cmdPath)
	if err := helm.Run(); err != nil {
		return errors.Wrapf(err, "cannot execute command: " + cmdPath)
	}

	return nil
}

func createHelmChartManifests(c *Config, appPath string) error {
	if err := helmCheckExecutable(); err != nil {
		return err
	}

	if output, outputStderr, err := helmRepoAdd(c.Spec.Helm.RepoName, c.Spec.Helm.RepoUrl); err != nil {
		log.Error("Command failed to run")
		log.Error(output)
		log.Error(outputStderr)
		return err
	}

	if output, outputStderr, err := helmRepoUpdate(); err != nil {
		log.Error("Command failed to run")
		log.Error(output)
		log.Error(outputStderr)
		return err
	}

	manifestDir, output, outputStderr, err := helmFetch(c.Spec.Helm.RepoName, c.Spec.Helm.ChartName, c.Spec.Helm.Version)
	if err != nil {
		log.Error("Command failed to run")
		log.Error(output)
		log.Error(outputStderr)
		return err
	}

	log.Info("Temporary helm chart dir: " + manifestDir)

	output, outputStderr, err = helmTemplate(appPath, manifestDir, c.Spec.Helm.ChartName, c.Spec.Namespace, c.Spec.Helm.ValuesFile)
	if err != nil {
		log.Error("Command failed to run")
		log.Error(output)
		log.Error(outputStderr)
		return err
	}

	log.Info("Removing temporary dir: " + manifestDir)
	if err := os.RemoveAll(manifestDir); err != nil {
		return err
	}

	return writeBaseHelmManifests(appPath, c.Spec.Helm.ChartName, output)
}

func helmRepoAdd(repoName, repoUrl string) ([]byte, []byte, error) {
	helmCmd := helmExecutablePath()
	log.Infof("Running: %s repo add %s %s", helmCmd, repoName, repoUrl)
	helm := exec.Command(helmCmd, "repo", "add", repoName, repoUrl)
	stderrBuf := new(bytes.Buffer)
	helm.Stderr = stderrBuf
	out, err := helm.Output()
	return out, stderrBuf.Bytes(), err
}

func helmRepoUpdate() ([]byte, []byte, error) {
	helmCmd := helmExecutablePath()
	log.Infof("Running: %s repo update", helmCmd)
	helm := exec.Command(helmCmd, "repo", "update")
	stderrBuf := new(bytes.Buffer)
	helm.Stderr = stderrBuf
	out, err := helm.Output()
	return out, stderrBuf.Bytes(), err
}

func helmFetch(repoName, chartName, version string) (string, []byte, []byte, error) {
	tmpDir, err := ioutil.TempDir("", "helm_fetch_tmp")
	if err != nil {
		return "", []byte{}, []byte{}, err
	}

	helmCmd := helmExecutablePath()
	repoSlashChart := fmt.Sprintf("%s/%s", repoName, chartName)
	log.Infof("Running: %s fetch %s --untar --version %s", helmCmd, repoSlashChart, version)

	helm := exec.Command(helmCmd, "fetch", repoSlashChart, "--untar", "--version", version)
	helm.Dir = tmpDir
	out, err := helm.Output()
	stderrBuf := new(bytes.Buffer)
	helm.Stderr = stderrBuf

	return tmpDir, out, stderrBuf.Bytes(), err
}

func helmTemplate(appPath, helmFetchDir, chartName, namespace, valuesFile string) ([]byte, []byte, error) {
	helmCmd := helmExecutablePath()
	log.Infof("Running: %s template --name-template %s --namespace %s -f %s --include-crds .", helmCmd, chartName, namespace, valuesFile)

	chartDirTmp := filepath.Join(helmFetchDir, chartName)
	exists, err := files.FileOrDirExists(chartDirTmp)
	if err != nil {
		return []byte{}, []byte{}, errors.Wrap(err, "Failed to check if helm chart tmp dir/chartName exists")
	}

	if !exists {
		return []byte{}, []byte{}, errors.New("Helm chart tmp dir/chartName does not exist")
	}

	valuesFileAbs := filepath.Join(appPath, valuesFile)

	helm := exec.Command(helmCmd, "template", "--name-template", chartName, "--namespace", namespace,
		"-f", valuesFileAbs, "--include-crds", ".")
	stderrBuff := new(bytes.Buffer)
	helm.Stderr = stderrBuff
	helm.Dir = chartDirTmp

	out, err := helm.Output()
	return out, stderrBuff.Bytes(), err
}

func writeBaseHelmManifests(appPath, chartName string, manifests []byte) error {
	path := filepath.Join(appPath, "base", "manifests", chartName+".yaml")
	return ioutil.WriteFile(path, manifests, os.FileMode(0640))
}