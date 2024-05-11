package client

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/eynopv/lac/internal/request"
	"github.com/eynopv/lac/internal/result"
)

type Client struct {
	Timeout int
}

func NewClient(timeout int) Client {
	return Client{
		Timeout: timeout,
	}
}

func (c *Client) Do(r *request.Request) (*result.Result, error) {
	var (
		request *http.Request
		err     error
	)

	if len(r.Body) != 0 {
		bodyReader := bytes.NewReader(r.Body)
		request, err = http.NewRequest(r.Method, r.Path, bodyReader)
	} else {
		request, err = http.NewRequest(r.Method, r.Path, nil)
	}

	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		request.Header.Set(key, value)
	}

	start := time.Now()
	client := http.Client{Timeout: time.Duration(c.Timeout) * time.Second}
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

	result, err := result.NewResult(
		elapsedTime,
		res.Request.URL.String(),
		res.Status,
		res.StatusCode,
		res.Header,
		body,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
