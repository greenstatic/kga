package files

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

func CreteDirs(pathToDirs []string, basePath string) error {
	for _, path := range pathToDirs {
		p := path
		if !filepath.IsAbs(p) && basePath != "" {
			p = filepath.Join(basePath, path)
		}
		if err := os.MkdirAll(p, os.FileMode(0775)); err != nil {
			return errors.Wrapf(err, "failed to create dir: %s", path)
		}
	}
	return nil
}

func RemoveDirIfExists(dir string) error {
	exists, err := FileOrDirExists(dir)
	if err != nil {
		return nil
	}

	if exists {
		return os.RemoveAll(dir)
	}

	return nil
}

func EnumerateAllFiles(path string) ([]string, error) {
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