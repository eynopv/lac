package request

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/eynopv/lac/internal/param"
	"github.com/eynopv/lac/internal/result"
	"github.com/eynopv/lac/internal/utils"
)

type Request struct {
	Method  string `json:"method" yaml:"method"`
	Path    string `json:"path" yaml:"path"`
	Body    json.RawMessage
	Headers map[string]string `json:"headers" yaml:"headers"`
}

func NewRequestV2(
	method string,
	path string,
	body []byte,
	headers map[string]string,
) Request {
	return Request{
		Method:  method,
		Path:    path,
		Body:    body,
		Headers: headers,
	}
}

func (r *Request) ResolveParameters(variables map[string]string) {
	r.Path = param.Param(r.Path).Resolve(variables)
	for key, value := range r.Headers {
		r.Headers[key] = param.Param(value).Resolve(variables)
	}
	if len(r.Body) != 0 {
		stringBody := param.Param(string(r.Body)).Resolve(variables)
		r.Body = []byte(stringBody)
	}
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

	requestData.Path = param.Param(requestData.Path).Resolve(variables)

	var request *http.Request

	if len(requestData.Body) != 0 {
		stringBody := param.Param(string(requestData.Body)).Resolve(variables)
		bodyReader := strings.NewReader(stringBody)
		request, err = http.NewRequest(requestData.Method, requestData.Path, bodyReader)
	} else {
		request, err = http.NewRequest(requestData.Method, requestData.Path, nil)
	}

	if err != nil {
		return nil, err
	}

	finalHeaders := utils.CombineMaps(headers, requestData.Headers)
	for key, value := range finalHeaders {
		request.Header.Set(key, param.Param(value).Resolve(variables))
	}

	return request, nil
}

func LoadRequest(itemPath string) (*Request, error) {
	var request Request
	if err := utils.LoadItem(itemPath, &request); err != nil {
		return nil, err
	}
	return &request, nil
}

func DoRequest(request *http.Request, timeout int) (*result.Result, error) {
	start := time.Now()
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	res, err := client.Do(request)
	elapsedTime := time.Since(start)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result, err := result.NewResult(
		elapsedTime,
		res.Status,
		res.StatusCode,
		res.Header,
		body,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
