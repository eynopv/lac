package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
	"time"

	"github.com/eynopv/gorcli/internal/utils"
)

type Request struct {
	Method  string `json:"method" yaml:"method"`
	Path    string `json:"path" yaml:"path"`
	Body    json.RawMessage
	Headers map[string]string `json:"headers" yaml:"headers"`
}

type Result struct {
	Status      string
	StatusCode  int
	ElapsedTime time.Duration
	Body        map[string]interface{}
	Headers     http.Header
}

func (r *Result) Print(showHeaders bool) {
	fmt.Println("Status:", r.Status)
	fmt.Println("Elapsed Time:", r.ElapsedTime)

	if showHeaders {
		utils.PrintPrettyJson(r.Headers)
	}

	if r.Body != nil {
		utils.PrintPrettyJson(r.Body)
	}
}

func LoadRequest(itemPath string) (*Request, error) {
	var request Request

	err := utils.LoadItem(itemPath, &request)

	if err != nil {
		fmt.Printf("Failed to parse request %s: %v\n", itemPath, err)
		return nil, err
	}

	return &request, nil
}

func NewRequest(
	name string,
	headers map[string]string,
	variables map[string]string,
) (*http.Request, error) {
	var requestData *Request
	var err error

	if requestData, err = LoadRequest(name); err != nil {
		return nil, err
	}

	requestData.Path = ParseParam(requestData.Path, variables)

	var request *http.Request

	if len(requestData.Body) != 0 {
		stringBody := ParseParam(string(requestData.Body), variables)
		bodyReader := strings.NewReader(stringBody)
		request, err = http.NewRequest(requestData.Method, requestData.Path, bodyReader)
	} else {
		request, err = http.NewRequest(requestData.Method, requestData.Path, nil)
	}

	if err != nil {
		return nil, err
	}

	finalHeaders := map[string]string{}

	if headers != nil {
		maps.Copy(finalHeaders, headers)
	}

	if requestData.Headers != nil {
		maps.Copy(finalHeaders, requestData.Headers)
	}

	for key, value := range finalHeaders {
		request.Header.Set(key, ParseParam(value, variables))
	}

	return request, nil
}

func DoRequest(request *http.Request, timeout int) (*Result, error) {
	result := Result{}

	start := time.Now()
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	res, err := client.Do(request)
	elapsed := time.Since(start)

	if err != nil {
		return nil, err
	}

	result.ElapsedTime = elapsed

	result.Status = res.Status
	result.StatusCode = res.StatusCode
	result.Headers = res.Header

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if len(body) > 0 {
		var responseData map[string]interface{}
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			fmt.Println("Failed to parse JSON:", err)
		}

		result.Body = responseData
	}

	return &result, nil
}
