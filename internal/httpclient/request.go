package httpclient

import (
	"bytes"
	"net/http"
	"net/url"
)

func BuildRequest(req Request) (*http.Request, error) {
	u, err := url.Parse(req.Url)
	if err != nil {
		return nil, err
	}

	var body *bytes.Reader
	if req.Body != "" {
		body = bytes.NewReader([]byte(req.Body))
	} else {
		body = bytes.NewReader(nil)
	}

	httpReq, err := http.NewRequest(req.Method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	return httpReq, nil
}
