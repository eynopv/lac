package request

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/eynopv/lac/pkg/param"
	"github.com/eynopv/lac/pkg/result"
	"github.com/eynopv/lac/pkg/utils"
)

type Request struct {
	Method  string `json:"method" yaml:"method"`
	Path    string `json:"path" yaml:"path"`
	Body    json.RawMessage
	Headers map[string]string `json:"headers" yaml:"headers"`
}

func LoadRequest(itemPath string) (*Request, error) {
	var request Request
	if err := utils.LoadItem(itemPath, &request); err != nil {
		return nil, err
	}
	return &request, nil
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
