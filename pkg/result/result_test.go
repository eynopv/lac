package result

import (
	"net/http"
	"testing"
	"time"
)

func TestNewResultValidInputs(t *testing.T) {
	result, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, nil, "HTTP/2.0", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Path != "/path" {
		t.Errorf("expected path '/path', got '%v'", result.Path)
	}
	if result.Status != "200 OK" {
		t.Errorf("expected status '200 OK', got '%v'", result.Status)
	}
	if result.StatusCode != 200 {
		t.Errorf("expected status code '200', got '%v'", result.StatusCode)
	}
	if result.Protocol != "HTTP/2.0" {
		t.Errorf("expected protocol 'HTTP/2.0', got '%v'", result.Protocol)
	}
	if result.Body != nil {
		t.Errorf("expected body to be nil, got %v", result.Body)
	}
	if result.Text != "" {
		t.Errorf(`expected text to be empty , got '%v'`, result.Text)
	}
}

func TestNewResultJsonContent(t *testing.T) {
	headers := http.Header{}
	headers.Set("content-type", "application/json")
	bodyRaw := []byte(`{"key":"value"}`)
	result, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, headers, "HTTP/2.0", bodyRaw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Body == nil || result.Body["key"] != "value" {
		t.Errorf("expected body with key 'key' and value 'value', got %v", result.Body)
	}
	if result.Text != "" {
		t.Errorf(`expected text to be empty , got '%v'`, result.Text)
	}
}

func TestNewResultTextContent(t *testing.T) {
	headers := http.Header{}
	headers.Set("content-type", "text/plain")
	bodyRaw := []byte("plain text content")
	result, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, headers, "HTTP/2.0", bodyRaw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Body != nil {
		t.Errorf("expected body to be nil, got %v", result.Body)
	}
	if result.Text != "plain text content" {
		t.Errorf(`expected text to be 'plain text content' , got '%v'`, result.Text)
	}
}

func TestNewResultEmptyBody(t *testing.T) {
	headers := http.Header{}
	headers.Set("content-type", "application/json")
	result, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, headers, "HTTP/2.0", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Body != nil {
		t.Errorf("expected body to be nil, got %v", result.Body)
	}
	if result.Text != "" {
		t.Errorf(`expected text to be empty , got '%v'`, result.Text)
	}
}

func TestNewResultInvalidJson(t *testing.T) {
	headers := http.Header{}
	headers.Set("content-type", "application/json")
	bodyRaw := []byte(`{invalid json}`)
	_, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, headers, "HTTP/2.0", bodyRaw)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
