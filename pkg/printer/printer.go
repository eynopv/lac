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
		req := *res.Response.Request
		s := ""
		s += p.formatter.RequestLine(*res.RequestLine())
		s += p.formatter.Headers(req.Header)
		sections = append(sections, s)
	}

	if p.config.PrintRequestBody {
		s := ""
		if requestJson := res.RequestJson(); requestJson != nil {
			s += fmt.Sprintf("%v\n", p.formatter.Json(requestJson))
		} else if requestText := res.RequestText(); requestText != "" {
			s += fmt.Sprintf("%v\n", requestText)
		}

		sections = append(sections, s)
	}

	if p.config.PrintResponseHeaders {
		s := ""
		s += p.formatter.StatusLine(*res.StatusLine())
		s += p.formatter.Headers(res.Response.Header)
		sections = append(sections, s)
	}

	if p.config.PrintResponseBody {
		s := ""
		if responseJson := res.Json(); responseJson != nil {
			s += fmt.Sprintf("%v\n", p.formatter.Json(responseJson))
		} else if responseText := res.Text(); responseText != "" {
			s += fmt.Sprintf("%v\n", responseText)
		}

		sections = append(sections, s)
	}

	fmt.Fprint(destination, strings.Join(sections, "\n"))
}
