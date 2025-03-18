package http_method

import (
	"net/http"
	"testing"

	"github.com/eynopv/lac/internal/assert"
)

func TestNormalizeHttpMethod(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		value string
	}{
		{
			name:  "unknown method",
			want:  "UNKNOWN",
			value: "this_is_invalid_method",
		},
		{
			name:  "get method",
			want:  http.MethodGet,
			value: "get",
		},
		{
			name:  "post method",
			want:  http.MethodPost,
			value: "post",
		},
		{
			name:  "put method",
			want:  http.MethodPut,
			value: "put",
		},
		{
			name:  "patch method",
			want:  http.MethodPatch,
			value: "patch",
		},
		{
			name:  "delete method",
			want:  http.MethodDelete,
			value: "delete",
		},
		{
			name:  "case insensitive",
			want:  http.MethodGet,
			value: "gEt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converted := NormalizeHttpMethod(tt.value)
			assert.Equal(t, converted, tt.want)
		})
	}
}
