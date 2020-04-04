package generate

import (
	"bytes"
	"fmt"
	"github.com/greenstatic/kga/internal/layout"
	"github.com/greenstatic/kga/pkg/config"
	"github.com/greenstatic/kga/pkg/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func helmExecutablePath() string {
	envHelm := os.Getenv("HELM")
	if envHelm != "" {
		return envHelm
	}

	return "helm" // assume that it is in the users path
}

func helmCheckExecutable() {
	cmdPath := helmExecutablePath()
	helm := exec.Command(cmdPath)
	if err := helm.Run(); err != nil {
		log.Error(err)
		log.Fatalf("Cannot execute command: %s", cmdPath)
	}
}

func CreateHelmChartManifests(appPath, namespace string, spec *config.HelmSpec) {
	helmCheckExecutable()

	if output, outputStderr, err := helmRepoAdd(spec.RepoName, spec.RepoUrl); err != nil {
		log.Error("Command failed to run")
		log.Error(output)
		log.Error(outputStderr)
		log.Fatal(err)
	}

	if output, outputStderr, err := helmRepoUpdate(); err != nil {
		log.Error("Command failed to run")
		log.Error(output)
		log.Error(outputStderr)
		log.Fatal(err)
	}

	manifestDir, output, outputStderr, err := helmFetch(spec.RepoName, spec.ChartName, spec.Version)
	if err != nil {
		log.Error("Command failed to run")
		log.Error(output)
		log.Error(outputStderr)
		log.Fatal(err)
	}

	log.Info("Temporary helm chart dir: " + manifestDir)

	output, outputStderr, err = helmTemplate(appPath, manifestDir, spec.ChartName, namespace, spec.ValuesFile)
	if err != nil {
		log.Error("Command failed to run")
		log.Error(output)
		log.Error(outputStderr)
		log.Fatal(err)
	}

	log.Info("Removing temporary directory: " + manifestDir)
	_ = os.RemoveAll(manifestDir) // cleanup tmp dir

	if err := layout.CreateBaseHelmManifests(appPath, spec.ChartName, output); err != nil {
		log.Fatal(err)
	}
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
	log.Infof("Running: %s template --name-template %s --namespace %s -f %s .", helmCmd, chartName, namespace, valuesFile)

	chartDirTmp := filepath.Join(helmFetchDir, chartName)
	exists, err := layout.FileOrDirExists(chartDirTmp)
	if err != nil {
		return []byte{}, []byte{}, errors.Wrap(err, "Failed to check if helm chart tmp dir/chartName exists")
	}

	if !exists {
		return []byte{}, []byte{}, errors.New("Helm chart tmp dir/chartName does not exist")
	}

	valuesFileAbs := filepath.Join(appPath, valuesFile)

	helm := exec.Command(helmCmd, "template", "--name-template", chartName, "--namespace", namespace,
		"-f", valuesFileAbs, ".")
	stderrBuff := new(bytes.Buffer)
	helm.Stderr = stderrBuff
	helm.Dir = chartDirTmp

	out, err := helm.Output()
	return out, stderrBuff.Bytes(), err
}
