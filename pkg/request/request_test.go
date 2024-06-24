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
}

func TestNewRequest(t *testing.T) {
	data := RequestData{
		Path:    "https://example.com",
		Method:  "get",
		Body:    []byte{},
		Headers: map[string]string{},
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
