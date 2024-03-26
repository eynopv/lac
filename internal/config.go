package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Loaded  bool
	Headers map[string]string `json:"headers"`
}

func LoadConfig() (*Config, error) {
	var config Config

	configPath := "./gorcli.config.json"
	filePath, err := FullPath(configPath)
	if err != nil {
		fmt.Println("Failed to load config file:", err)
		return nil, err
	}
	configExists := FileExists(*filePath)

	if configExists {
		content, err := os.ReadFile("./gorcli.config.json")
		if err != nil {
			fmt.Println("Failed to load config file:", err)
			return nil, err
		}

		err = json.Unmarshal(content, &config)
		if err != nil {
			fmt.Println("Failed to load config file:", err)
			return nil, err
		}

		for key, value := range config.Headers {
			param := Param{Name: key, Value: value}
			config.Headers[key] = param.ParseValue()
		}

		return &config, nil
	}

	config = Config{
		Headers: map[string]string{},
	}

	return &config, nil
}
