package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type GorRequest struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   Body   `json:"body"`
}

type Body map[string]interface{}

type Config struct {
	Headers map[string]string `json:"headers"`
}

func main() {
	config, err := loadConfig()
	if err != nil {
		return
	}

	fmt.Println(config.Headers)

	requestName := os.Args[1]
	request, err := buildRequest(requestName)
	if err != nil {
		fmt.Println("Unable to make request:", err)
		return
	}

	for key, value := range config.Headers {
		fmt.Printf("Setting Header %s = %s", key, value)
		request.Header.Set(key, value)
	}

	client := http.Client{}
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer res.Body.Close()

	fmt.Println("Status:", res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed to read resonse body:", err)
		return
	}

	if len(body) == 0 {
		fmt.Println("Response body is empty")
		return
	}

	prettyHeaders, err := json.MarshalIndent(res.Header, "", " ")
	fmt.Println(string(prettyHeaders))

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
	}

	prettyJSON, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

func loadRequest(name string) (*GorRequest, error) {
	var request GorRequest

	filePath := fmt.Sprintf("./collections/%s.json", name)
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to load file %s: %v", filePath, err)
		return nil, err
	}

	err = json.Unmarshal(content, &request)
	if err != nil {
		fmt.Printf("Failed to load file %s: %v", filePath, err)
		return nil, err
	}

	return &request, nil
}

func buildRequest(name string) (*http.Request, error) {
	requestData, err := loadRequest(name)

	if err != nil {
		return nil, err
	}

	var body *strings.Reader
	body = nil
	if requestData.Body != nil {
		body, err = buildBody(requestData.Body)

		if err != nil {
			return nil, err
		}
	}

	var request *http.Request
	if body != nil {
		request, err = http.NewRequest(requestData.Method, requestData.Path, body)
	} else {
		request, err = http.NewRequest(requestData.Method, requestData.Path, nil)
	}

	if err != nil {
		return nil, err
	}

	return request, nil
}

func loadConfig() (*Config, error) {
	var config Config

	configPath := "./gorcli.json"
	configExists := FileExists(configPath)

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

		return &config, nil
	}

	config = Config{
		Headers: map[string]string{},
	}

	return &config, nil
}

func buildBody(body Body) (*strings.Reader, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(string(bodyBytes)), nil
}
