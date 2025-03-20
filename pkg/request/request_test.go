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
}

func TestNewRequest(t *testing.T) {
	data := RequestData{
		Path:    "https://example.com",
		Method:  "get",
		Body:    []byte{},
		Headers: map[string]StringOrStringList{},
	}
	request := NewRequest(data)
	assert.Equal(t, request.Path, data.Path)
	assert.Equal(t, request.Method, http.MethodGet)
	assert.DeepEqual(t, request.Body, data.Body)
	assert.DeepEqual(t, request.Headers, data.Headers)
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
  }
}
`
	err := json.Unmarshal([]byte(data), &requestData)

	assert.NoError(t, err)
	assert.Equal(t, len(requestData.Headers), 2)
	assert.Equal(t, requestData.Headers["Accept"][0], "text/plain")
	assert.Equal(t, requestData.Headers["Accept"][1], "application/json")
}
