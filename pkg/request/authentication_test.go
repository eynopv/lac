package request

import (
	"net/http"
	"testing"

	"github.com/eynopv/lac/internal/assert"
	"github.com/eynopv/lac/internal/errorsx"
)

func TestNewAuth(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		template := Template(`
    auth:
      type: unknown
    `)

		auth, err := NewAuth(&template)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errorsx.ErrAuthUnknown.Error())
		assert.Nil(t, auth)
	})

	t.Run("json", func(t *testing.T) {
		template := Template(`
    {
      "auth": {
        "type": "unknown",
      }
    }
    `)

		auth, err := NewAuth(&template)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errorsx.ErrAuthUnknown.Error())
		assert.Nil(t, auth)
	})

	t.Run("no auth", func(t *testing.T) {
		template := Template(`
    {
      "hello": "world"
    }
    `)

		auth, err := NewAuth(&template)

		assert.NoError(t, err)
		assert.Nil(t, auth)
	})

	t.Run("basic", func(t *testing.T) {
		template := Template(`
    auth:
      type: basic
      username: hello
      password: world
    `)

		auth, err := NewAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.GetType(), Basic)
	})

	t.Run("bearer", func(t *testing.T) {
		template := Template(`
    auth:
      type: bearer
      token: helloworld
    `)

		auth, err := NewAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.GetType(), Bearer)
	})
}

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
    "username": "hello"
  }
}
`
		template := Template(data)
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
		template := Template(data)
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
		template := Template(data)
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

	request, err := http.NewRequest(http.MethodGet, "", nil)

	assert.NotNil(t, request)
	assert.NoError(t, err)

	basic.Apply(request)

	username, password, ok := request.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, basic.Username, username)
	assert.Equal(t, basic.Password, password)
}

func TestNewBearerAuth(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		template := Template(`
    auth:
      token: helloworld
    `)

		auth, err := NewBearerAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.Token, "helloworld")
	})

	t.Run("json", func(t *testing.T) {
		template := Template(`
    {
      "auth": {
        "token": "helloworld"
      }
    }
    `)

		auth, err := NewBearerAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.Token, "helloworld")
	})

	t.Run("invalid", func(t *testing.T) {
		template := Template("this is invalid template")

		auth, err := NewBearerAuth(&template)

		assert.Error(t, err)
		assert.Nil(t, auth)
	})

	t.Run("template without auth", func(t *testing.T) {
		template := Template(`
    {
      "hello": "world"
    }
    `)

		auth, err := NewBearerAuth(&template)

		assert.NoError(t, err)
		assert.True(t, auth == nil)
	})
}

func TestBearerAuthApply(t *testing.T) {
	bearer := BearerAuth{
		Token: "helloworld",
	}

	request, err := http.NewRequest(http.MethodGet, "", nil)

	assert.NotNil(t, request)
	assert.NoError(t, err)

	bearer.Apply(request)

	assert.Equal(t, request.Header.Get("Authorization"), "Bearer helloworld")
}
