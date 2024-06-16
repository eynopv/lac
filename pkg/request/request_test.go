package request

import (
	"strings"
	"testing"
)

func TestRequest(t *testing.T) {
	t.Run("resolve parameters", func(t *testing.T) {
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
	})
}
