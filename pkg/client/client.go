package client

import (
	"io"
	"net/http"
	"time"

	"github.com/eynopv/lac/pkg/authentication"
	"github.com/eynopv/lac/pkg/printer"
	"github.com/eynopv/lac/pkg/request"
	"github.com/eynopv/lac/pkg/result"
)

type Client struct {
	timeout         int
	followRedirects bool
	PrinterConfig   printer.PrinterConfig
}

type ClientConfig struct {
	Timeout       int
	NoRedirects   bool
	PrinterConfig printer.PrinterConfig
}

func NewClient(config *ClientConfig) *Client {
	return &Client{
		timeout:         config.Timeout,
		followRedirects: !config.NoRedirects,
		PrinterConfig:   config.PrinterConfig,
	}
}

func (c *Client) Do(r *request.Request, auth authentication.Auth) (*result.Result, error) {
	request, err := r.ToHttpRequest()
	if err != nil {
		return nil, err
	}

	if auth != nil {
		auth.Apply(request)
	}

	client := c.ToHttpClient()
	start := time.Now()
	res, err := client.Do(request)
	elapsedTime := time.Since(start)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return &result.Result{
		Response:     res,
		ResponseBody: body,
		RequestBody:  r.Body,
		Metadata: result.Metadata{
			ElapsedTime: elapsedTime,
		},
	}, nil
}

func (c *Client) ToHttpClient() *http.Client {
	client := http.Client{Timeout: time.Duration(c.timeout) * time.Second}
	if !c.followRedirects {
		client.CheckRedirect = NoRedirects
	}

	return &client
}

func NoRedirects(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
