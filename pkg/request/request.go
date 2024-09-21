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
	Method    string `json:"method" yaml:"method"`
	Path      string `json:"path" yaml:"path"`
	Body      json.RawMessage
	Headers   map[string]string `json:"headers" yaml:"headers"`
	Variables map[string]string `json:"variables" yaml:"variables"`
}

type Request struct {
	Method    string
	Path      string
	Body      []byte
	Headers   map[string]string
	Variables map[string]string
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
		Method:    utils.StringToHttpMethod(data.Method),
		Path:      data.Path,
		Body:      data.Body,
		Headers:   data.Headers,
		Variables: data.Variables,
	}
}

func (r *Request) ResolveParameters(variables map[string]string) {
	r.resolvePathParameters(variables)
	r.resolveHeaderParameters(variables)
	r.resolveBodyParameters(variables)
}

func (r *Request) resolvePathParameters(variables map[string]string) {
	initialPath := r.Path
	r.Path = param.Param(r.Path).Resolve(variables)
	if r.Path == initialPath {
		r.Path = param.Param(r.Path).Resolve(r.Variables)
	}
}

func (r *Request) resolveHeaderParameters(variables map[string]string) {
	for key, value := range r.Headers {
		initialHeaderValue := r.Headers[strings.ToLower(key)]
		r.Headers[strings.ToLower(key)] = param.Param(value).Resolve(variables)
		if initialHeaderValue == r.Headers[strings.ToLower(key)] {
			r.Headers[strings.ToLower(key)] = param.Param(value).Resolve(r.Variables)
		}
	}
}

func (r *Request) resolveBodyParameters(variables map[string]string) {
	if len(r.Body) != 0 {
		initialBody := string(r.Body)
		stringBody := param.Param(string(r.Body)).Resolve(variables)
		if initialBody == stringBody {
			stringBody = param.Param(string(r.Body)).Resolve(r.Variables)
		}
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

func (r *Request) ToHttpRequest() (*http.Request, error) {
	var (
		request *http.Request
		err     error
	)

	if len(r.Body) != 0 {
		bodyStr := strings.Trim(strings.TrimSpace(string(r.Body)), `"`)
		bodyReader := strings.NewReader(bodyStr)
		request, err = http.NewRequest(r.Method, r.Path, bodyReader)
	} else {
		request, err = http.NewRequest(r.Method, r.Path, nil)
	}

	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		request.Header.Set(key, value)
	}

	return request, nil
}
