package client

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	clientConfig := ClientConfig{
		Timeout:     30,
		NoRedirects: true,
	}
	client := NewClient(&clientConfig)
	if client.timeout != 30 {
		t.Fatalf("expected timeout %d to be %d", client.timeout, clientConfig.Timeout)
	}
	if client.followRedirects != false {
		t.Fatalf("expected followeRedirects %v to be %v", client.followRedirects, clientConfig.Timeout)
	}
}

func TestToHttpClient(t *testing.T) {
	t.Run("has default follow redirects policy", func(t *testing.T) {
		client := NewClient(&ClientConfig{
			Timeout:     15,
			NoRedirects: false,
		})
		httpClient := client.ToHttpClient()
		if httpClient.CheckRedirect != nil {
			t.Fatalf("expected CheckRedirect to be nil")
		}
	})

	t.Run("no follow redirects", func(t *testing.T) {
		client := NewClient(&ClientConfig{
			Timeout:     15,
			NoRedirects: true,
		})
		httpClient := client.ToHttpClient()
		if httpClient.CheckRedirect == nil {
			t.Fatalf("expected CheckRedirect to be set")
		}
	})
}

func TestNotRedirectsReturnsLastResponse(t *testing.T) {
	err := NoRedirects(nil, nil)
	if err != http.ErrUseLastResponse {
		t.Fatalf("expected error to be ErrUseLastResponse")
	}
}
