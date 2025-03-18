package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/eynopv/lac/pkg/http_method"
	"github.com/eynopv/lac/pkg/param"
	"github.com/eynopv/lac/pkg/utils"
	yaml "gopkg.in/yaml.v3"
)

type RequestData struct {
	Method    string                        `json:"method" yaml:"method"`
	Path      string                        `json:"path" yaml:"path"`
	Body      ByteBody                      `json:"body" yaml:"body"`
	Headers   map[string]StringOrStringList `json:"headers" yaml:"headers"`
	Variables map[string]string             `json:"variables" yaml:"variables"`
}

type Request struct {
	Method    string
	Path      string
	Body      []byte
	Headers   map[string]StringOrStringList
	Variables map[string]string
}

type ByteBody []byte

func (b *ByteBody) UnmarshalJSON(data []byte) error {
	*b = ByteBody(data)
	return nil
}

func (b *ByteBody) UnmarshalYAML(value *yaml.Node) error {
	var raw interface{}
	if err := value.Decode(&raw); err != nil {
		return err
	}

	jsonData, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	*b = ByteBody(jsonData)

	return nil
}

type StringOrStringList []string

func (s *StringOrStringList) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = []string{single}
		return nil
	}

	var multiple []string
	if err := json.Unmarshal(data, &multiple); err == nil {
		*s = multiple
		return nil
	}

	return fmt.Errorf("invalid format")
}

func (s *StringOrStringList) UnmarshalYAML(value *yaml.Node) error {
	var raw interface{}
	if err := value.Decode(&raw); err != nil {
		return err
	}

	jsonData, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	return s.UnmarshalJSON(jsonData)
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
	method := http.MethodGet
	if data.Method != "" {
		method = http_method.NormalizeHttpMethod(data.Method)
	}

	return Request{
		Method:    method,
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
		var resolvedValues []string

		for _, v := range value {
			resolved := param.Param(v).Resolve(variables)
			if resolved == v {
				resolved = param.Param(v).Resolve(r.Variables)
			}

			resolvedValues = append(resolvedValues, resolved)
		}

		r.Headers[strings.ToLower(key)] = resolvedValues
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
		request.Header[key] = value
	}

	return request, nil
}
