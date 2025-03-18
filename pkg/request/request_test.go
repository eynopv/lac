package request

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/eynopv/lac/internal/assert"
	"gopkg.in/yaml.v3"
)

func TestNewRequestDefaults(t *testing.T) {
	request := NewRequest(RequestData{})
	assert.Equal(t, request.Path, "")
	assert.Equal(t, request.Method, "GET")
	assert.DeepEqual(t, request.Body, nil)
	assert.DeepEqual(t, request.Headers, nil)
	assert.DeepEqual(t, request.Variables, nil)
}

func TestNewRequest(t *testing.T) {
	data := RequestData{
		Path:      "https://example.com",
		Method:    "get",
		Body:      []byte{},
		Headers:   map[string]StringOrStringList{},
		Variables: map[string]string{},
	}
	request := NewRequest(data)
	assert.Equal(t, request.Path, data.Path)
	assert.Equal(t, request.Method, http.MethodGet)
	assert.DeepEqual(t, request.Body, data.Body)
	assert.DeepEqual(t, request.Headers, data.Headers)
	assert.DeepEqual(t, request.Variables, data.Variables)
}

func TestRequestResolveParameters(t *testing.T) {
	bodyString := "{\"field\":\"${value}\"}"

	headers := map[string]StringOrStringList{
		"x-api-key": []string{"${API_KEY}"},
	}

	variables := map[string]string{
		"id":      "resolvedId",
		"API_KEY": "resolvedApiKey",
		"value":   "resolvedValue",
	}

	req := Request{
		Method:  "POST",
		Path:    "example.com/api/${id}",
		Body:    []byte(bodyString),
		Headers: headers,
	}

	req.ResolveParameters(variables)

	assert.Equal(t, req.Path, "example.com/api/resolvedId")
	assert.Equal(t, req.Headers["x-api-key"][0], variables["API_KEY"])
	assert.StringContains(t, string(req.Body), variables["value"])
}

func TestResolveHeaderParameterFromRequestVariable(t *testing.T) {
	headers := map[string]StringOrStringList{
		"var": []string{"${REQUEST_VARIABLE}"},
	}

	req := Request{
		Headers: headers,
		Variables: map[string]string{
			"REQUEST_VARIABLE": "resolved",
		},
	}
	req.ResolveParameters(nil)

	assert.Equal(t, req.Headers["var"][0], "resolved")
}

func TestResolvePathParameterFromRequestVariable(t *testing.T) {
	initialPath := "example.com/api/${id}"
	req := Request{
		Path: initialPath,
		Variables: map[string]string{
			"id": "resolved",
		},
	}
	req.ResolveParameters(nil)
	assert.NotEqual(t, req.Path, initialPath)
}

func TestResolveBodyParameterFromRequestVariabes(t *testing.T) {
	initialBody := "{\"field\":\"${value}\"}"
	req := Request{
		Body: []byte(initialBody),
		Variables: map[string]string{
			"value": "resolved",
		},
	}
	req.ResolveParameters(nil)
	assert.NotEqual(t, string(req.Body), initialBody)
}

func TestUnmarshalYaml(t *testing.T) {
	var requestData RequestData

	data := `
path: ${host}/post
method: POST
headers:
  Content-Type: application/json
  Accept:
    - text/plain
    - application/json
body:
  key: value
variables:
  host: https://example.com
  `
	err := yaml.Unmarshal([]byte(data), &requestData)

	assert.NoError(t, err)
	assert.Equal(t, len(requestData.Headers), 2)
	assert.Equal(t, requestData.Headers["Accept"][0], "text/plain")
	assert.Equal(t, requestData.Headers["Accept"][1], "application/json")
}

func TestUnmarshalJson(t *testing.T) {
	var requestData RequestData

	data := `
{
  "path": "${host}/post",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json",
    "Accept": ["text/plain", "application/json"]
  },
  "body": {
    "key": "value"
  },
  "variables": {
    "host": "https://example.com"
  }
}
`
	err := json.Unmarshal([]byte(data), &requestData)

	assert.NoError(t, err)
	assert.Equal(t, len(requestData.Headers), 2)
	assert.Equal(t, requestData.Headers["Accept"][0], "text/plain")
	assert.Equal(t, requestData.Headers["Accept"][1], "application/json")
}
