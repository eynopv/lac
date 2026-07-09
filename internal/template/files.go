package template

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var ErrFileExists = errors.New("file already exists")
var ErrFileDoesNotExists = errors.New("file already exists")

func Exists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

func Load(f string) (*Template, error) {
	if !Exists(f) {
		return nil, ErrFileDoesNotExists
	}

	data, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var t Template

	err = json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func Save(filename string, t Template) error {
	if Exists(filename) {
		return ErrFileExists
	}

	dir := filepath.Dir(filename)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
