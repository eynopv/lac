package result

import (
	"encoding/json"
	"net/http"
	"time"
)

type Result struct {
	Response     *http.Response
	ResponseBody Body
	RequestBody  Body
	Metadata     Metadata
}

type Metadata struct {
	ElapsedTime time.Duration
}

type StatusLine struct {
	Protocol string        `json:"protocol"`
	Status   string        `json:"status"`
	Time     time.Duration `json:"time"`
}

type RequestLine struct {
	Protocol string `json:"protocol"`
	Url      string `json:"url"`
	Method   string `json:"method"`
}

type Body []byte

func (r Result) StatusLine() *StatusLine {
	return &StatusLine{
		Protocol: r.Response.Proto,
		Status:   r.Response.Status,
		Time:     r.Metadata.ElapsedTime,
	}
}

func (r Result) ResponseLine() *StatusLine {
	return r.StatusLine()
}

func (r Result) RequestLine() *RequestLine {
	return &RequestLine{
		Protocol: r.Response.Request.Proto,
		Url:      r.Response.Request.URL.String(),
		Method:   r.Response.Request.Method,
	}
}

func (b Body) Data() any {
	if body := b.Json(); body != nil {
		return body
	}

	return b.Text()
}

func (b Body) Json() map[string]any {
	if len(b) == 0 {
		return nil
	}

	var data map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		return nil
	}

	return data
}

func (b Body) Text() string {
	return string(b)
}
