package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/eynopv/lac/pkg/result"
)

type Includes struct {
	ResponseBody    bool
	ResponseHeaders bool
	ResponseMeta    bool
	RequestBody     bool
	RequestHeaders  bool
	RequestMeta     bool
}

type Outputter struct {
	includes    Includes
	destination io.Writer
	formatter   Formatter
}

func NewOutputter(config PrinterConfig) *Outputter {
	includes := Includes{
		ResponseBody:    config.PrintRequestBody,
		ResponseHeaders: config.PrintResponseHeaders,
		ResponseMeta:    config.PrintResponseMeta,
		RequestBody:     config.PrintRequestBody,
		RequestHeaders:  config.PrintRequestHeaders,
		RequestMeta:     config.PrintRequestMeta,
	}

	var formatter Formatter
	if IsTerminal(int(os.Stdout.Fd())) {
		formatter = &PrettyFormatter{
			includes: includes,
		}
	} else {
		formatter = &JsonFormatter{
			includes: includes,
		}
	}

	return &Outputter{
		includes:    includes,
		destination: os.Stdout,
		formatter:   formatter,
	}
}

func (o *Outputter) Write(res *result.Result) error {
	formatted, err := o.formatter.Format(res)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(o.destination, formatted)

	return err
}
