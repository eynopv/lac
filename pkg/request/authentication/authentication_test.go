package authentication

import (
	"testing"

	"github.com/eynopv/lac/internal/assert"
	"github.com/eynopv/lac/internal/errorsx"
	"github.com/eynopv/lac/pkg/request"
)

func TestNewAuth(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		template := request.Template(`
    auth:
      type: unknown
    `)

		auth, err := NewAuth(&template)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errorsx.ErrAuthUnknown.Error())
		assert.Nil(t, auth)
	})

	t.Run("json", func(t *testing.T) {
		template := request.Template(`
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
		template := request.Template(`
    {
      "hello": "world"
    }
    `)

		auth, err := NewAuth(&template)

		assert.NoError(t, err)
		assert.Nil(t, auth)
	})

	t.Run("basic", func(t *testing.T) {
		template := request.Template(`
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
		template := request.Template(`
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
