package app

import (
	"github.com/greenstatic/kga/pkg/log"
	"os"
	"os/exec"
)

func kustomizeExecutablePath() string {
	env := os.Getenv("KUSTOMIZE")
	if env != "" {
		return env
	}
	return "kustomize" // assume that is is in the users path
}

func kustomizeBuildSuccseeds(path string) error {
	cmd := kustomizeExecutablePath()
	log.Infof("Checking if kustomize can build: %s build %s", cmd, path)
	k := exec.Command(cmd, "build", path)
	return k.Run()
}
