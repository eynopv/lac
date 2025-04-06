package printer

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/eynopv/lac/internal/assert"
	"github.com/eynopv/lac/pkg/result"
)

func Test_formatStatusLine(t *testing.T) {
	tests := []struct {
		name     string
		line     result.StatusLine
		useColor bool
		want     string
	}{
		{
			name: "non-colorized 200",
			line: result.StatusLine{
				Protocol: "HTTP/1.1",
				Status:   "200 OK",
				Time:     300 * time.Millisecond,
			},
			useColor: false,
			want:     "HTTP/1.1 200 OK [300ms]\n",
		},
		{
			name: "colorized 200 fast response",
			line: result.StatusLine{
				Protocol: "HTTP/2",
				Status:   "200 OK",
				Time:     300 * time.Millisecond,
			},
			useColor: true,
			want: fmt.Sprintf(
				"%v %v [%v]\n",
				"HTTP/2",
				Colorize("200 OK", ColorGreen),
				Colorize("300ms", ColorReset),
			),
		},
		{
			name: "colorized 301 medium response",
			line: result.StatusLine{
				Protocol: "HTTP/1.1",
				Status:   "301 Moved Permanently",
				Time:     800 * time.Millisecond,
			},
			useColor: true,
			want: fmt.Sprintf(
				"%v %v [%v]\n",
				"HTTP/1.1",
				Colorize("301 Moved Permanently", ColorCyan),
				Colorize("800ms", ColorYellow),
			),
		},
		{
			name: "colorized 500 slow response",
			line: result.StatusLine{
				Protocol: "HTTP/1.1",
				Status:   "500 Internal Server Error",
				Time:     1500 * time.Millisecond,
			},
			useColor: true,
			want: fmt.Sprintf(
				"%v %v [%v]\n",
				"HTTP/1.1",
				Colorize("500 Internal Server Error", ColorRed),
				Colorize("1.5s", ColorRed),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var formatter interface {
				StatusLine(result.StatusLine) string
			}
			if tt.useColor {
				formatter = ColorFormatter{}
			} else {
				formatter = PlainFormatter{}
			}
			got := formatter.StatusLine(tt.line)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_formatRequestLine(t *testing.T) {
	testUrl := "https://example.com/dogs"
	tests := []struct {
		name     string
		line     result.RequestLine
		useColor bool
		want     string
	}{
		{
			name: "non-colorized GET",
			line: result.RequestLine{
				Method:   "GET",
				Url:      testUrl,
				Protocol: "HTTP/1.1",
			},
			useColor: false,
			want:     fmt.Sprintf("GET %v HTTP/1.1\n", testUrl),
		},
		{
			name: "colorized GET",
			line: result.RequestLine{
				Method:   "GET",
				Url:      testUrl,
				Protocol: "HTTP/1.1",
			},
			useColor: true,
			want: fmt.Sprintf(
				"%v %v %v\n",
				Colorize("GET", ColorGreen),
				testUrl,
				"HTTP/1.1",
			),
		},
		{
			name: "colorized POST",
			line: result.RequestLine{
				Method:   "POST",
				Url:      testUrl,
				Protocol: "HTTP/2",
			},
			useColor: true,
			want: fmt.Sprintf(
				"%v %v %v\n",
				Colorize("POST", ColorYellow),
				testUrl,
				"HTTP/2",
			),
		},
		{
			name: "colorized DELETE",
			line: result.RequestLine{
				Method:   "DELETE",
				Url:      testUrl,
				Protocol: "HTTP/1.1",
			},
			useColor: true,
			want: fmt.Sprintf(
				"%v %v %v\n",
				Colorize("DELETE", ColorRed),
				testUrl,
				"HTTP/1.1",
			),
		},
		{
			name: "colorized default",
			line: result.RequestLine{
				Method:   "OPTIONS",
				Url:      testUrl,
				Protocol: "HTTP/1.0",
			},
			useColor: true,
			want: fmt.Sprintf(
				"%v %v %v\n",
				Colorize("OPTIONS", ColorMagenta),
				testUrl,
				"HTTP/1.0",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var formatter interface {
				RequestLine(result.RequestLine) string
			}
			if tt.useColor {
				formatter = ColorFormatter{}
			} else {
				formatter = PlainFormatter{}
			}
			got := formatter.RequestLine(tt.line)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_formatJson(t *testing.T) {
	tests := []struct {
		name string
		json map[string]interface{}
		want string
	}{
		{
			name: "simple",
			json: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			want: "{\n \"age\": 30,\n \"name\": \"Alice\"\n}",
		},
		{
			name: "nested",
			json: map[string]interface{}{
				"user": map[string]interface{}{
					"id":   1,
					"name": "Bob",
				},
			},
			want: "{\n \"user\": {\n  \"id\": 1,\n  \"name\": \"Bob\"\n }\n}",
		},
		{
			name: "empty",
			json: map[string]interface{}{},
			want: "{}",
		},
		{
			name: "unmarshalable value (channel)",
			json: map[string]interface{}{
				"invalid": make(chan int),
			},
			want: "unable to parse json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := PlainFormatter{}
			got := formatter.Json(tt.json)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_formatHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  http.Header
		useColor bool
		want     string
	}{
		{
			name:     "no headers",
			headers:  http.Header{},
			useColor: false,
			want:     "",
		},
		{
			name: "single header non-colorized",
			headers: http.Header{
				"Content-Type": {"application/json"},
			},
			useColor: false,
			want:     "Content-Type: application/json\n",
		},
		{
			name: "single header colorized",
			headers: http.Header{
				"Content-Type": {"application/json"},
			},
			useColor: true,
			want:     fmt.Sprintf("%v: %v\n", Cyan("Content-Type"), "application/json"),
		},
		{
			name: "multiple headers",
			headers: http.Header{
				"Content-Type": {"application/json"},
				"X-Custom":     {"value1", "value2"},
			},
			useColor: false,
			want:     "Content-Type: application/json\nX-Custom: value1, value2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var formatter interface {
				Headers(http.Header) string
			}
			if tt.useColor {
				formatter = ColorFormatter{}
			} else {
				formatter = PlainFormatter{}
			}
			got := formatter.Headers(tt.headers)
			assert.Equal(t, tt.want, got)
		})
	}
}
