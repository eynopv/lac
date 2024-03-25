package main

import (
	"os"
	"path/filepath"
)

func FileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil
}

func FullPath(fpath string) (*string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(cwd, fpath)
	return &fullPath, nil
}
