package request

import (
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

func TestTemplateParse(t *testing.T) {
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
