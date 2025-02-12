package client

import (
	"net/http"
	"testing"

	"github.com/eynopv/lac/internal/assert"
)

func TestNewClient(t *testing.T) {
	clientConfig := ClientConfig{
		Timeout:     30,
		NoRedirects: true,
	}
	client := NewClient(&clientConfig)
	assert.Equal(t, client.timeout, 30)
	assert.Equal(t, client.followRedirects, false)
}

func TestToHttpClient(t *testing.T) {
	t.Run("has default follow redirects policy", func(t *testing.T) {
		client := NewClient(&ClientConfig{
			Timeout:     15,
			NoRedirects: false,
		})
		httpClient := client.ToHttpClient()
		assert.DeepEqual(t, httpClient.CheckRedirect, nil)
	})

	t.Run("no follow redirects", func(t *testing.T) {
		client := NewClient(&ClientConfig{
			Timeout:     15,
			NoRedirects: true,
		})
		httpClient := client.ToHttpClient()
		assert.NotNil(t, httpClient.CheckRedirect)
	})
}

func TestNotRedirectsReturnsLastResponse(t *testing.T) {
	err := NoRedirects(nil, nil)
	assert.Equal(t, err, http.ErrUseLastResponse)
}
