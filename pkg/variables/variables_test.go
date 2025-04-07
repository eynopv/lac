package variables

import (
	"testing"

	"github.com/eynopv/lac/internal/assert"
	yaml "gopkg.in/yaml.v3"
)

func TestVariables_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "boolean",
			data:    []byte(`{"boolean": true}`),
			wantErr: false,
		},
		{
			name:    "integer",
			data:    []byte(`{"integer": 3}`),
			wantErr: false,
		},
		{
			name:    "float",
			data:    []byte(`{"float": 3.21}`),
			wantErr: false,
		},
		{
			name:    "string",
			data:    []byte(`{"string": "hello, world"}`),
			wantErr: false,
		},
		{
			name:    "list of same type",
			data:    []byte(`{"list": [1,2,3]}`),
			wantErr: true,
		},
		{
			name:    "list of different type",
			data:    []byte(`{"list": [1,"hello"]}`),
			wantErr: true,
		},
		{
			name:    "nested",
			data:    []byte(`{"parent": { "child": "hello" }}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Variables
			err := v.UnmarshalJSON(tt.data)

			if tt.wantErr {
				assert.ErrorContains(t, err, "unsupported variables type")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVariables_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		value   *yaml.Node
		wantErr bool
	}{
		{
			name:    "boolean",
			value:   yamlNode(`boolean: true`),
			wantErr: false,
		},
		{
			name:    "integer",
			value:   yamlNode(`integer: 3`),
			wantErr: false,
		},
		{
			name:    "float",
			value:   yamlNode(`float: 3.21`),
			wantErr: false,
		},
		{
			name:    "string",
			value:   yamlNode(`string: "hello, world"`),
			wantErr: false,
		},
		{
			name:    "list of same type",
			value:   yamlNode(`list: [1, 2, 3]`),
			wantErr: true,
		},
		{
			name:    "list of different type",
			value:   yamlNode(`list: [1, "hello"]`),
			wantErr: true,
		},
		{
			name:    "nested",
			value:   yamlNode(`parent: {child: "hello"}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Variables
			err := v.UnmarshalYAML(tt.value)

			if tt.wantErr {
				assert.ErrorContains(t, err, "unsupported variables type")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func yamlNode(s string) *yaml.Node {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(s), &node); err != nil {
		return nil
	}

	return node.Content[0]
}
