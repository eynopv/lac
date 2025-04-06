package printer

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/eynopv/lac/internal/assert"
	"github.com/eynopv/lac/pkg/result"
)

func TestNewPrinter(t *testing.T) {
	tests := []struct {
		name        string
		isTerminal  bool
		expectColor bool
	}{
		{"TerminalOutput", true, true},
		{"NonTerminalOutput", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalIsTerminal := IsTerminal
			defer func() { IsTerminal = originalIsTerminal }()

			IsTerminal = func(fd int) bool { return tt.isTerminal }

			p := NewPrinter(PrinterConfig{})

			assert.Equal(t, p.formatter.colored, tt.isTerminal)
		})
	}
}

func TestPrinter_Print(t *testing.T) {
	var buf bytes.Buffer

	res := result.Result{
		Response: &http.Response{
			Status: "200 OK",
			Proto:  "HTTP/1.1",
			Request: &http.Request{
				Method: "GET",
				Proto:  "HTTP/1.1",
				URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/"},
				Header: http.Header{"Req-H": []string{"v"}},
			},
			Header: http.Header{"Res-H": []string{"x"}},
		},
		RequestBody:  []byte(`request body`),
		ResponseBody: []byte(`response body`),
		Metadata:     result.Metadata{ElapsedTime: 123 * time.Millisecond},
	}

	cases := []struct {
		name   string
		config PrinterConfig
		expect []string
	}{
		{
			"OnlyResponseBody",
			PrinterConfig{PrintResponseBody: true},
			[]string{`response body`},
		},
		{
			"RequestHeadersAndBody",
			PrinterConfig{PrintRequestHeaders: true, PrintRequestBody: true},
			[]string{"GET https://example.com/", `request body`, "Req-H"},
		},
		{
			"AllSections",
			PrinterConfig{
				PrintRequestHeaders:  true,
				PrintRequestBody:     true,
				PrintResponseHeaders: true,
				PrintResponseBody:    true,
			},
			[]string{
				"GET https://example.com/ HTTP/1.1",
				"Req-H",
				"request body",
				"HTTP/1.1 200 OK [123ms]",
				"Res-H",
				"response body",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			p := Printer{
				config:      tt.config,
				destination: &buf,
				formatter:   Formatter{},
			}

			p.Print(&res)

			out := buf.String()
			for _, expect := range tt.expect {
				if !strings.Contains(out, expect) {
					t.Errorf("expected output to contain %q: %v", expect, out)
				}
			}
		})
	}
}
