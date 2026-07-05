package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/eynopv/lac/internal/template"
	"github.com/eynopv/lac/internal/utils"
)

var ErrFileExists = errors.New("file already exists")

func New(filename string) error {
	if utils.Exists(filename) {
		return ErrFileExists
	}

	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	r := template.Schema{}
	jsonBytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
