package printer

import (
	"io"

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

func (o *Outputter) Write(res *result.Result) error {
	
}
