package request

import (
	"os"
	"testing"

	"github.com/eynopv/lac/internal/assert"
)

func TestNewTemplate(t *testing.T) {
	oldFileLoader := fileLoader
	fileLoader = func(file string) (*[]byte, error) {
		data := []byte("Template Content")
		return &data, nil
	}

	defer func() { fileLoader = oldFileLoader }()

	template, err := NewTemplate("example.template")

	assert.NoError(t, err)
	assert.Equal(t, *template, "Template Content")
}

func TestTemplate_Parse(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		data := `
path: ${host}/post
method: POST
headers:
  Content-Type: application/json
  Accept:
    - text/plain
    - application/json
body:
  key: value
  `

		template := Template(data)
		request, err := template.Parse()

		assert.NoError(t, err)
		assert.NotNil(t, request)
	})

	t.Run("json", func(t *testing.T) {
		data := `
{
  "path": "${host}/post",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json",
    "Accept": ["text/plain", "application/json"]
  },
  "body": {
    "key": "value"
  }
}
`

		template := Template(data)
		request, err := template.Parse()

		assert.NoError(t, err)
		assert.NotNil(t, request)
	})

	t.Run("invalid", func(t *testing.T) {
		data := "this is invalid template"

		template := Template(data)
		request, err := template.Parse()

		assert.Error(t, err)
		assert.Nil(t, request)
	})
}

func TestTemplate_Interpolate(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    Template
	}{
		{
			name:    "json",
			content: `{"string": "${string}", "number": "${number}", "env_string": "${ENV_VAR}"}`,
			want:    Template(`{"string": "hello, world", "number": 7, "env_string": "dlrow, olleh"}`),
		},
		{
			name: "yaml",
			content: `
			string: ${string}
			number: ${number}
			env_string: ${ENV_VAR}
			`,
			want: Template(`
			string: hello, world
			number: 7
			env_string: dlrow, olleh
			`),
		},
		{
			name:    "form body",
			content: `{"body": "value=${string}&envValue=${ENV_VAR}"}`,
			want:    Template(`{"body": "value=hello, world&envValue=dlrow, olleh"}`),
		},
	}

	vars := map[string]any{
		"string": "hello, world",
		"number": 7,
	}

	os.Setenv("ENV_VAR", "dlrow, olleh")
	defer os.Unsetenv("ENV_VAR")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := Template(tt.content)

			got := template.Interpolate(vars, false)

			assert.Equal(t, *got, tt.want)
		})
	}
}
