package request

import (
	"testing"

	"github.com/eynopv/lac/internal/assert"
)

func TestNewBasicAuth(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		data := `
auth:
  username: hello
  password: world
  `

		template := Template(data)
		auth := NewBasicAuth(&template)

		assert.NotNil(t, auth)
	})

	t.Run("json", func(t *testing.T) {
		data := `
{
  "auth": {
    "username": "hello",
    "password": "world"
  }
}
`
		template := Template(data)
		auth := NewBasicAuth(&template)

		assert.NotNil(t, auth)
	})

	t.Run("invalid", func(t *testing.T) {
		data := "this is invalid template"

		template := Template(data)
		auth := NewBasicAuth(&template)

		assert.Nil(t, auth)
	})
}
