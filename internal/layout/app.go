package layout

import (
	"github.com/greenstatic/kga/pkg/log"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func PathIsKgaApp(appPath string) (bool, error) {
	return FileOrDirExists(filepath.Join(appPath, "kga.yaml"))
}

func CreateAppDir(directory string, appName string) error {
	appPath := filepath.Join(directory, appName)
	if err := appPathExists(appPath); err != nil {
		return err
	}

	return os.Mkdir(appPath, os.FileMode(0775))
}

func RemoveAppBaseDir(appPath string) error {
	basePath := filepath.Join(appPath, "base")
	exists, err := FileOrDirExists(basePath)
	if err != nil {
		log.Error(err)
		log.Fatal("Failed to remote app base dir")
	}

	if exists {
		return os.RemoveAll(basePath)
	}

	return nil
}

func appPathExists(appPath string) error {
	exists, err := FileOrDirExists(appPath)
	if err != nil {
		return errors.Wrap(err, "Failed to check if app directory exists")
	}

	if exists {
		return errors.New("App directory already exists, quitting")
	}

	return nil
}

// Source: https://stackoverflow.com/a/10510783
func FileOrDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
