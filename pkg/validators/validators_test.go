package validators

import (
	"testing"

	"github.com/eynopv/lac/internal/assert"
)

func TestIsJson(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{"valid object", []byte(`{"name": "John", "surname": "Doe"}`), true},
		{"valid array", []byte(`[1,2,3]`), true},
		{"valid empty object", []byte(`{}`), true},
		{"valid empty array", []byte(`[]`), true},
		{"invalid trailing comma", []byte(`{"name": "John",}`), false},
		{"invalid whitespace", []byte(` `), false},
		{"invalid empty", []byte(``), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsJson(tt.input)
			assert.Equal(t, result, tt.expected)
		})
	}
}
