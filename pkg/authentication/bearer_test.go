package authentication

import (
	"net/http"
	"testing"

	"github.com/eynopv/lac/internal/assert"
	"github.com/eynopv/lac/pkg/template"
)

func TestNewBearerAuth(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		template := template.Template(`
    auth:
      token: helloworld
    `)

		auth, err := NewBearerAuth(&template)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, auth.Token, "helloworld")
	})

	t.Run("json", func(t *testing.T) {
		template := template.Template(`
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
		template := template.Template("this is invalid template")

		auth, err := NewBearerAuth(&template)

		assert.Error(t, err)
		assert.Nil(t, auth)
	})

	t.Run("template without auth", func(t *testing.T) {
		template := template.Template(`
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

	template, err := http.NewRequest(http.MethodGet, "", nil)

	assert.NotNil(t, template)
	assert.NoError(t, err)

	bearer.Apply(template)

	assert.Equal(t, template.Header.Get("Authorization"), "Bearer helloworld")
}
