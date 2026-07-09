package httpclient

import (
	"io"
	"net/http"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Send(req Request) (*Response, error) {
	httpReq, err := BuildRequest(req)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:     resp.StatusCode,
		StatusText: resp.Status,
		Headers:    resp.Header,
		Body:       string(bodyBytes),
	}, nil
}
