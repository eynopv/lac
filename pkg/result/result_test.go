package result

import (
	"net/http"
	"testing"
	"time"

	"github.com/eynopv/lac/internal/assert"
)

func TestNewResultValidInputs(t *testing.T) {
	result, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, nil, "HTTP/2.0", nil)

	assert.NoError(t, err)
	assert.Equal(t, result.Path, "/path")
	assert.Equal(t, result.Status, "200 OK")
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Protocol, "HTTP/2.0")
	assert.Equal(t, result.Text, "")
	assert.DeepEqual(t, result.Body, nil)
}

func TestNewResultJsonContent(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	bodyRaw := []byte(`{"key":"value"}`)
	result, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, headers, "HTTP/2.0", bodyRaw)

	assert.NoError(t, err)
	assert.Equal(t, result.Body["key"], "value")
	assert.Equal(t, result.Text, "")
}

func TestNewResultTextContent(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "text/plain")

	bodyRaw := []byte("plain text content")
	result, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, headers, "HTTP/2.0", bodyRaw)

	assert.NoError(t, err)
	assert.DeepEqual(t, result.Body, nil)
	assert.Equal(t, result.Text, "plain text content")
}

func TestNewResultEmptyBody(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	result, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, headers, "HTTP/2.0", nil)

	assert.NoError(t, err)
	assert.DeepEqual(t, result.Body, nil)
	assert.Equal(t, result.Text, "")
}

func TestNewResultInvalidJson(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	bodyRaw := []byte(`{invalid json}`)
	_, err := NewResult(100*time.Millisecond, "/path", "200 OK", 200, headers, "HTTP/2.0", bodyRaw)

	assert.Error(t, err)
}
