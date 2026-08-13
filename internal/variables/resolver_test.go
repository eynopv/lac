package variables_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eynopv/lac/internal/variables"
)

func TestResolver_Resolve(t *testing.T) {
	tests := []struct {
		name  string
		vars  map[string]any
		input string
		want  string
	}{
		{
			name:  "string",
			vars:  map[string]any{"placeholder": "Resolver"},
			input: "Hello, ${placeholder}!",
			want:  "Hello, Resolver!",
		},
		{
			name:  "quoted string",
			vars:  map[string]any{"placeholder": "Resolver"},
			input: "Hello, \"${placeholder}\"!",
			want:  "Hello, \"Resolver\"!",
		},
		{
			name:  "boolean",
			vars:  map[string]any{"placeholder": true},
			input: "{ \"enabled\": ${placeholder} }",
			want:  "{ \"enabled\": true }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := variables.NewResolver(tt.vars)
			got := r.Resolve(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
