package printer

import (
	"io"
	"os"

	"golang.org/x/term"
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
	return Printer{
		config:      config,
		destination: os.Stdout,
	}
}
