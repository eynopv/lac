package authentication

import (
	"net/http"
	"testing"

	"github.com/eynopv/lac/internal/assert"
	"github.com/eynopv/lac/pkg/template"
)

func TestNewApiAuth(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		template := template.Template(`
    auth:
      header: x-api-key
      key: helloworld
    `)

		auth, err := NewApiAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.Header, "x-api-key")
		assert.Equal(t, auth.Key, "helloworld")
	})

	t.Run("json", func(t *testing.T) {
		template := template.Template(`
    {
      "auth": {
        "header": "x-api-key",
        "key": "helloworld"
      }
    }
    `)

		auth, err := NewApiAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.Header, "x-api-key")
		assert.Equal(t, auth.Key, "helloworld")
	})

	t.Run("invalid", func(t *testing.T) {
		template := template.Template("this is invalid template")

		auth, err := NewApiAuth(&template)

		assert.Error(t, err)
		assert.Nil(t, auth)
	})

	t.Run("template without auth", func(t *testing.T) {
		template := template.Template(`
    {
      "hello": "world"
    }
    `)

		auth, err := NewApiAuth(&template)

		assert.NoError(t, err)
		assert.True(t, auth == nil)
	})
}

func TestApiAuthApply(t *testing.T) {
	apiAuth := ApiAuth{
		Header: "x-api-key",
		Key:    "helloworld",
	}

	template, err := http.NewRequest(http.MethodGet, "", nil)

	assert.NotNil(t, template)
	assert.NoError(t, err)

	apiAuth.Apply(template)

	assert.Equal(t, template.Header.Get("X-Api-Key"), "helloworld")
}
