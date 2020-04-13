package app

import (
	"fmt"
	"github.com/greenstatic/kga/pkg/files"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func PathAbsolute(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(pwd, path)
}

type PathExistsError string

func (s PathExistsError) Error() string {
	return fmt.Sprintf("app path: %s exists", string(s))
}

// Checks whether the appPath (which can be relative or absolute) is valid to initialized
func ValidAppPathToInit(appPath string) error {
	exists, err := files.FileOrDirExists(appPath)
	if err != nil {
		return errors.Wrap(err, "failed to check if dir exists")
	}

	if exists {
		return PathExistsError(appPath)
	}

	return nil
}
