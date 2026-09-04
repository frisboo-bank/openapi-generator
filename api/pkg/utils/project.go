package utils

import (
	"os"
	"path/filepath"
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/syserrors"
)

var (
	ErrReadingDirectory = syserrors.New("error reading directory")
	ErrNoProjectFound   = syserrors.New("error failed to find any go.mod file")
)

func GetProjectRootWorkingDirectory() (string, error) {
	pwdD, _ := os.Getwd()
	dir, err := searchRootDirectory(pwdD)
	if err != nil {
		return "", err
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func searchRootDirectory(baseDir string) (string, error) {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return "", syserrors.WrapWith(ErrReadingDirectory, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if strings.EqualFold(filename, "go.mod") {
			return baseDir, nil
		}
	}

	parentDir := filepath.Dir(baseDir)
	if parentDir == baseDir {
		return "", syserrors.WrapWith(ErrNoProjectFound, err)
	}

	return searchRootDirectory(parentDir)
}
