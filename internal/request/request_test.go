package request

import (
	"strings"
	"testing"
)

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

	req := NewRequestV2("POST", "example.com/api/${id}", []byte(bodyString), headers)
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
