package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadItem(item string, dst interface{}) error {
	var (
		err  error
		data *[]byte
	)

	if data, err = LoadFile(item); err != nil {
		return err
	}

	if strings.HasSuffix(item, ".json") {
		return json.Unmarshal(*data, dst)
	}

	if strings.HasSuffix(item, ".yaml") || strings.HasSuffix(item, ".yml") {
		return yaml.Unmarshal(*data, dst)
	}

	return fmt.Errorf("Not supported file: %v", item)
}

func LoadFile(file string) (*[]byte, error) {
	var (
		filePath string
		err      error
		content  []byte
	)

	if filePath, err = GetFullPath(file); err != nil {
		return nil, err
	}

	if fileExists := PathExists(filePath); !fileExists {
		return nil, err
	}

	if content, err = os.ReadFile(filePath); err != nil {
		return nil, err
	}

	return &content, nil
}

func PathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func GetFullPath(filePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(cwd, filePath)

	return fullPath, nil
}
