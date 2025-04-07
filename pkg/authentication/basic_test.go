package authentication

import (
	"net/http"
	"testing"

	"github.com/eynopv/lac/internal/assert"
	"github.com/eynopv/lac/internal/errorsx"
	"github.com/eynopv/lac/pkg/template"
)

func TestNewBasicAuth(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		data := `
auth:
  username: hello
  password: world
  `

		template := template.Template(data)
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
		template := template.Template(data)
		auth, err := NewBasicAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.Username, "hello")
		assert.Equal(t, auth.Password, "world")
	})

	t.Run("invalid", func(t *testing.T) {
		data := "this is invalid template"

		template := template.Template(data)
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

		template := template.Template(data)
		auth, err := NewBasicAuth(&template)

		assert.NoError(t, err)
		assert.True(t, auth == nil)
	})

	t.Run("only username", func(t *testing.T) {
		data := `
{
  "auth": {
    "username": "hello"
  }
}
`
		template := template.Template(data)
		auth, err := NewBasicAuth(&template)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errorsx.ErrBasicAuthInvalid.Error())
		assert.Nil(t, auth)
	})

	t.Run("only password", func(t *testing.T) {
		data := `
{
  "auth": {
    "password": "hello"
  }
}
`
		template := template.Template(data)
		auth, err := NewBasicAuth(&template)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errorsx.ErrBasicAuthInvalid.Error())
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
		template := template.Template(data)
		auth, err := NewBasicAuth(&template)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errorsx.ErrBasicAuthInvalid.Error())
		assert.Nil(t, auth)
	})
}

func TestBasicAuthApply(t *testing.T) {
	basic := BasicAuth{
		Username: "hello",
		Password: "world",
	}

	template, err := http.NewRequest(http.MethodGet, "", nil)

	assert.NotNil(t, template)
	assert.NoError(t, err)

	basic.Apply(template)

	username, password, ok := template.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, basic.Username, username)
	assert.Equal(t, basic.Password, password)
}
