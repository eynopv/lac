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
	PrintResponseMeta    bool
	PrintRequestBody     bool
	PrintRequestHeaders  bool
	PrintRequestMeta     bool
}

type Printer struct {
	config      PrinterConfig
	destination io.Writer
	formatter   Formatter
}

func NewPrinter(config PrinterConfig) Printer {
	formatter := Formatter{
		colored: IsTerminal(int(os.Stdout.Fd())),
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

	if p.config.PrintRequestMeta {
		return p.formatter.RequestLine(*res.RequestLine()) + p.formatter.Headers(req.Header)
	}

	return p.formatter.Headers(req.Header)
}

func (p *Printer) printResponseHeaders(res *result.Result) string {
	if p.config.PrintResponseMeta {
		return p.formatter.StatusLine(*res.StatusLine()) + p.formatter.Headers(res.Response.Header)
	}

	return p.formatter.Headers(res.Response.Header)
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
