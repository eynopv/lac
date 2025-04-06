package request

import (
	"encoding/json"
	"fmt"
	yaml "gopkg.in/yaml.v3"
	"net/http"
	"strings"

	"github.com/eynopv/lac/pkg/http_method"
	"github.com/eynopv/lac/pkg/utils"
)

type RequestData struct {
	Method  string                        `json:"method" yaml:"method"`
	Path    string                        `json:"path" yaml:"path"`
	Body    ByteBody                      `json:"body" yaml:"body"`
	Headers map[string]StringOrStringList `json:"headers" yaml:"headers"`
}

type Request struct {
	Method  string
	Path    string
	Body    []byte
	Headers map[string]StringOrStringList
}

type ByteBody []byte

func (b *ByteBody) UnmarshalJSON(data []byte) error {
	*b = ByteBody(data)
	return nil
}

func (b *ByteBody) UnmarshalYAML(value *yaml.Node) error {
	var raw any
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
	var raw any
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

	var headers map[string]StringOrStringList

	if data.Headers != nil {
		headers = map[string]StringOrStringList{}
		for key, value := range data.Headers {
			headers[http.CanonicalHeaderKey(key)] = value
		}
	}

	return Request{
		Method:  method,
		Path:    data.Path,
		Body:    data.Body,
		Headers: headers,
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
