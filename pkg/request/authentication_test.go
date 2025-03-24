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
		auth, err := NewBasicAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.Username, "hello")
		assert.Equal(t, auth.Password, "world")
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
		auth, err := NewBasicAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.Username, "hello")
		assert.Equal(t, auth.Password, "world")
	})

	t.Run("invalid", func(t *testing.T) {
		data := "this is invalid template"

		template := Template(data)
		auth, err := NewBasicAuth(&template)

		assert.Error(t, err)
		assert.Nil(t, auth)
	})

	t.Run("template without auth", func(t *testing.T) {
		data := `
{
  "hello": "world"
}
`

		template := Template(data)
		auth, err := NewBasicAuth(&template)

		assert.NoError(t, err)
		assert.True(t, auth == nil)
	})

	t.Run("only username", func(t *testing.T) {
		data := `
{
  "auth": {
    "username": "hello",
  }
}
`
		template := Template(data)
		auth, err := NewBasicAuth(&template)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), ErrBasicAuthInvalid.Error())
		assert.Nil(t, auth)
	})

	t.Run("only password", func(t *testing.T) {
		data := `
{
  "auth": {
    "password": "hello",
  }
}
`
		template := Template(data)
		auth, err := NewBasicAuth(&template)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), ErrBasicAuthInvalid.Error())
		assert.Nil(t, auth)
	})

	t.Run("empty username and password", func(t *testing.T) {
		data := `
{
  "auth": {
    "username": "",
    "password": ""
  }
}
`
		template := Template(data)
		auth, err := NewBasicAuth(&template)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), ErrBasicAuthInvalid.Error())
		assert.Nil(t, auth)
	})
}
