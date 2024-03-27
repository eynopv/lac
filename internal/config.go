package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Headers   map[string]string `json:"headers"`
	Variables map[string]string `json:"variables"`
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

		// TODO: headers could use variables so maybe params should be parsed right before request
		for key, value := range config.Headers {
			param := Param{Name: key, Value: value}
			config.Headers[key] = param.ParseValue()
		}

		return &config, nil
	}

	config = Config{
		Headers:   map[string]string{},
		Variables: map[string]string{},
	}

	return &config, nil
}
