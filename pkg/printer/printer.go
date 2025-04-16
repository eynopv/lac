package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/eynopv/lac/pkg/result"
)

var IsTerminal = term.IsTerminal

type PrinterConfig struct {
	PrintResponseBody    bool
	PrintResponseHeaders bool
	PrintRequestBody     bool
	PrintRequestHeaders  bool
}

type Printer struct {
	config      PrinterConfig
	destination io.Writer
	formatter   Formatter
}

func NewPrinter(config PrinterConfig) Printer {
	formatter := Formatter{
		colored: true,
	}

	return Printer{
		config:      config,
		destination: os.Stdout,
		formatter:   formatter,
	}
}

func (p *Printer) Print(res *result.Result) {
	var output string

	if res.Response == nil {
		fmt.Fprint(p.destination, "No HTTP response available\n")
		return
	}

	if IsTerminal(int(os.Stdout.Fd())) {
		output = p.makeTerminalOutput(res)
	} else {
		output = p.makeNonTerminalOutput(res)
	}

	fmt.Fprint(p.destination, output)
}

func (p *Printer) makeTerminalOutput(res *result.Result) string {
	sections := []string{}

	if p.config.PrintRequestHeaders {
		sections = append(sections, p.printRequestHeaders(res))
	}

	if p.config.PrintRequestBody {
		sections = append(sections, p.printBody(&res.RequestBody))
	}

	if p.config.PrintResponseHeaders {
		sections = append(sections, p.printResponseHeaders(res))
	}

	if p.config.PrintResponseBody {
		sections = append(sections, p.printBody(&res.ResponseBody))
	}

	return strings.Join(sections, "\n")
}

func (p *Printer) makeNonTerminalOutput(res *result.Result) string {
	m := map[string]any{}

	if p.config.PrintRequestHeaders || p.config.PrintRequestBody {
		rm := map[string]any{}

		if p.config.PrintRequestHeaders {
			rm["headers"] = res.Response.Request.Header
		}

		if p.config.PrintRequestBody {
			if jsonBody := res.RequestBody.Json(); jsonBody != nil {
				rm["body"] = jsonBody
			} else if textBody := res.RequestBody.Text(); textBody != "" {
				rm["body"] = textBody
			}
		}

		m["request"] = rm
	}

	if p.config.PrintResponseHeaders || p.config.PrintResponseBody {
		rm := map[string]any{}

		if p.config.PrintResponseHeaders {
			rm["headers"] = res.Response.Header
		}

		if p.config.PrintResponseBody {
			if jsonBody := res.ResponseBody.Json(); jsonBody != nil {
				rm["body"] = jsonBody
			} else if textBody := res.ResponseBody.Text(); textBody != "" {
				rm["body"] = textBody
			}
		}

		m["response"] = rm
	}

	b, err := json.MarshalIndent(m, "", "  ")

	if err != nil {
		return err.Error()
	}

	return string(b)
}

func (p *Printer) printRequestHeaders(res *result.Result) string {
	req := *res.Response.Request
	return p.formatter.RequestLine(*res.RequestLine()) + p.formatter.Headers(req.Header)
}

func (p *Printer) printResponseHeaders(res *result.Result) string {
	return p.formatter.StatusLine(*res.StatusLine()) + p.formatter.Headers(res.Response.Header)
}

func (p *Printer) printBody(body *result.Body) string {
	if jsonBody := body.Json(); jsonBody != nil {
		return fmt.Sprintf("%v\n", p.formatter.Json(jsonBody))
	}

	if textBody := body.Text(); textBody != "" {
		return fmt.Sprintf("%v\n", textBody)
	}

	return ""
}
