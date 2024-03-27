package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Request struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   json.RawMessage
}

func LoadRequest(name string) (*Request, error) {
	var request Request

	filePath := fmt.Sprintf("./collections/%s.json", name)
	fullPath, err := FullPath(filePath)
	if err != nil {
		fmt.Printf("Failed to load file %s: %v\n", *fullPath, err)
		return nil, err
	}

	content, err := os.ReadFile(*fullPath)
	if err != nil {
		fmt.Printf("Failed to load file %s: %v\n", *fullPath, err)
		return nil, err
	}

	err = json.Unmarshal(content, &request)
	if err != nil {
		fmt.Printf("Failed to load file %s: %v\n", *fullPath, err)
		return nil, err
	}

	return &request, nil
}

func NewRequest(name string) (*http.Request, error) {
	requestData, err := LoadRequest(name)

	if err != nil {
		return nil, err
	}

	var request *http.Request
	if len(requestData.Body) != 0 {
		bodyReader := bytes.NewReader(requestData.Body)
		request, err = http.NewRequest(requestData.Method, ParseStringParam(requestData.Path), bodyReader)
	} else {
		request, err = http.NewRequest(requestData.Method, ParseStringParam(requestData.Path), nil)
	}

	if err != nil {
		return nil, err
	}

	return request, nil
}
