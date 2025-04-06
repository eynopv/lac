package printer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/eynopv/lac/pkg/result"
)

type PrinterConfig struct {
	PrintResponseBody    bool
	PrintResponseHeaders bool
	PrintRequestBody     bool
	PrintRequestHeaders  bool
}

type Printer struct {
	config    PrinterConfig
	formatter Formatter
}

var destination io.Writer = os.Stdout
var isTerminal = term.IsTerminal(int(os.Stdout.Fd()))

func NewPrinter(config PrinterConfig) Printer {
	formatter := NewFormatter(isTerminal)

	return Printer{
		config:    config,
		formatter: formatter,
	}
}

func (p *Printer) Print(res *result.Result) {
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

	fmt.Fprint(destination, strings.Join(sections, "\n"))
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
