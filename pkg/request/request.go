package request

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/eynopv/lac/pkg/param"
	"github.com/eynopv/lac/pkg/result"
	"github.com/eynopv/lac/pkg/utils"
)

type RequestData struct {
	Method  string `json:"method" yaml:"method"`
	Path    string `json:"path" yaml:"path"`
	Body    json.RawMessage
	Headers map[string]string `json:"headers" yaml:"headers"`
}

type Request struct {
	Method  string
	Path    string
	Body    []byte
	Headers map[string]string
}

func LoadRequest(itemPath string) (*Request, error) {
	var data RequestData
	if err := utils.LoadItem(itemPath, &data); err != nil {
		return nil, err
	}
	request := NewRequest(data)
	return &request, nil
}

func NewRequest(data RequestData) Request {
	return Request{
		Method:  utils.StringToHttpMethod(data.Method),
		Path:    data.Path,
		Body:    data.Body,
		Headers: data.Headers,
	}
}

func (r *Request) ResolveParameters(variables map[string]string) {
	r.Path = param.Param(r.Path).Resolve(variables)
	for key, value := range r.Headers {
		r.Headers[strings.ToLower(key)] = param.Param(value).Resolve(variables)
	}
	if len(r.Body) != 0 {
		stringBody := param.Param(string(r.Body)).Resolve(variables)
		r.Body = []byte(stringBody)
	}
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
		res.Request.URL.Path,
		res.Status,
		res.StatusCode,
		res.Header,
		res.Proto,
		body,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
