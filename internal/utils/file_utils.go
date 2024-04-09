package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadItem(item string, dst interface{}) error {
	var data *[]byte
	data = LoadFile(fmt.Sprintf("%s.json", item))
	if data != nil {
		err := json.Unmarshal(*data, dst)
		return err
	}

	data = LoadFile(fmt.Sprintf("%s.yaml", item))
	if data != nil {
		err := yaml.Unmarshal(*data, dst)
		return err
	}

	return errors.New("Not found")
}

func LoadFile(file string) *[]byte {
	var filePath string
	var err error

	if filePath, err = FullPath(file); err != nil {
		//fmt.Printf("Failed to resolve full path to file %s: %v\n", file, err)
		return nil
	}

	if fileExists := FileExists(filePath); !fileExists {
		//fmt.Printf("File %s does not exist\n", filePath)
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		//fmt.Printf("Failed to read file %s: %v\n", file, err)
		return nil
	}

	return &content
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func FullPath(filePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(cwd, filePath)
	return fullPath, nil
}
