package result

import (
	"encoding/json"
	"net/http"
	"time"
)

type Result struct {
	Response     *http.Response
	ResponseBody []byte
	RequestBody  []byte
	Metadata     Metadata
}

type Metadata struct {
	ElapsedTime time.Duration
}

func (r Result) StatusLine() *StatusLine {
	return &StatusLine{
		Protocol: r.Response.Proto,
		Status:   r.Response.Status,
		Time:     r.Metadata.ElapsedTime,
	}
}

func (r Result) RequestLine() *RequestLine {
	return &RequestLine{
		Protocol: r.Response.Request.Proto,
		Url:      r.Response.Request.URL.String(),
		Method:   r.Response.Request.Method,
	}
}

func (r Result) RequestJson() map[string]interface{} {
	body := r.RequestBody

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil
	}

	return data
}

func (r Result) RequestText() string {
	body := r.RequestBody
	return string(body)
}

func (r Result) Json() map[string]interface{} {
	if len(r.ResponseBody) == 0 {
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(r.ResponseBody, &data); err != nil {
		return nil
	}

	return data
}

func (r Result) Text() string {
	return string(r.ResponseBody)
}

type StatusLine struct {
	Protocol string
	Status   string
	Time     time.Duration
}

type RequestLine struct {
	Protocol string
	Url      string
	Method   string
}
