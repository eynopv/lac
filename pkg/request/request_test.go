package request

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestNewRequestDefaults(t *testing.T) {
	request := NewRequest(RequestData{})
	if request.Path != "" {
		t.Fatalf("expected path to be empty, got '%s'", request.Path)
	}
	if request.Method != "UNKNOWN METHOD" {
		t.Fatalf("expected method to be 'UNKNOWN METHOD', got '%s'", request.Method)
	}
	if request.Body != nil {
		t.Fatalf("expected body to be empty, got '%v'", request.Body)
	}
	if request.Headers != nil {
		t.Fatalf("expected headers to be empty, got '%v'", request.Headers)
	}
	if request.Variables != nil {
		t.Fatalf("expected variables to be empty, got '%v'", request.Variables)
	}
}

func TestNewRequest(t *testing.T) {
	data := RequestData{
		Path:      "https://example.com",
		Method:    "get",
		Body:      []byte{},
		Headers:   map[string]string{},
		Variables: map[string]string{},
	}
	request := NewRequest(data)
	if request.Path != data.Path {
		t.Fatalf("expected path to be '%s', got '%s'", data.Path, request.Path)
	}
	if request.Method != http.MethodGet {
		t.Fatalf("expected method to be '%s', got '%s'", http.MethodGet, request.Method)
	}
	if !reflect.DeepEqual(request.Body, []byte{}) {
		t.Fatalf("expected body to be empty slice, got '%v'", len(request.Body))
	}
	if !reflect.DeepEqual(request.Headers, map[string]string{}) {
		t.Fatalf("expected headers to be empty map, got '%v'", request.Headers)
	}
	if !reflect.DeepEqual(request.Variables, map[string]string{}) {
		t.Fatalf("expected variables to be empty map, got '%v'", request.Variables)
	}
}

func TestRequestResolveParameters(t *testing.T) {
	bodyString := "{\"field\":\"${value}\"}"

	headers := map[string]string{
		"x-api-key": "${API_KEY}",
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

	if req.Path != "example.com/api/resolvedId" {
		t.Fatalf("Path was not resolved: " + req.Path)
	}

	if req.Headers["x-api-key"] != "resolvedApiKey" {
		t.Fatalf("Headers were not resolved: " + req.Headers["x-api-key"])
	}

	if !strings.Contains(string(req.Body), "resolvedValue") {
		t.Fatalf("Body was not resolved: " + string(req.Body))
	}
}

func TestResolveHeaderParameterFromRequestVariable(t *testing.T) {
	headers := map[string]string{
		"var": "${REQUEST_VARIABLE}",
	}

	req := Request{
		Headers: headers,
		Variables: map[string]string{
			"REQUEST_VARIABLE": "resolved",
		},
	}
	req.ResolveParameters(nil)

	if req.Headers["var"] != "resolved" {
		t.Fatalf("headers parameter was not resolved: %v", req.Headers["var"])
	}
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
	if req.Path == initialPath {
		t.Fatalf("path was not resolved: %v", req.Path)
	}
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
	if string(req.Body) == initialBody {
		t.Fatalf("body was not resolved: %v", string(req.Body))
	}
}
