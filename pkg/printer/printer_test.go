package printer

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

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
			oldIsTerminal := isTerminal
			defer func() { isTerminal = oldIsTerminal }()

			isTerminal = tt.isTerminal

			p := NewPrinter(PrinterConfig{})

			switch f := p.formatter.(type) {
			case ColorFormatter:
				if !tt.expectColor {
					t.Errorf("expected PlainFormatter, got ColorFormatter")
				}
			case PlainFormatter:
				if tt.expectColor {
					t.Errorf("expected ColorFormatter, got PlainFormatter")
				}
			default:
				t.Errorf("unexpected formatter type: %T", f)
			}
		})
	}
}

func TestPrinter_Print(t *testing.T) {
	var buf bytes.Buffer

	oldDestination := destination
	defer func() { destination = oldDestination }()
	destination = &buf

	oldIsTerminal := isTerminal
	defer func() { isTerminal = oldIsTerminal }()
	isTerminal = false

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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			p := NewPrinter(tc.config)
			p.Print(&res)

			out := buf.String()
			for _, expect := range tc.expect {
				if !strings.Contains(out, expect) {
					t.Errorf("expected output to contain %q: %v", expect, out)
				}
			}
		})
	}
}
