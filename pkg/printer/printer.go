package printer

import (
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
	var formatter Formatter
	if IsTerminal(int(os.Stdout.Fd())) {
		formatter = ColorFormatter{}
	} else {
		formatter = PlainFormatter{}
	}

	return Printer{
		config:      config,
		destination: os.Stdout,
		formatter:   formatter,
	}
}

func (p *Printer) Print(res *result.Result) {
	if res.Response == nil {
		fmt.Fprint(p.destination, "No HTTP response available\n")
		return
	}

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

	fmt.Fprint(p.destination, strings.Join(sections, "\n"))
}

func (p *Printer) printRequestHeaders(res *result.Result) string {
	req := *res.Response.Request
	s := ""
	s += p.formatter.RequestLine(*res.RequestLine())
	s += p.formatter.Headers(req.Header)

	return s
}

func (p *Printer) printResponseHeaders(res *result.Result) string {
	s := ""
	s += p.formatter.StatusLine(*res.StatusLine())
	s += p.formatter.Headers(res.Response.Header)

	return s
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
