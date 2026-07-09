package app

import (
	"errors"

	"github.com/eynopv/lac/internal/httpclient"
	"github.com/eynopv/lac/internal/output"
	"github.com/eynopv/lac/internal/template"
)

var ErrFileDoesNotExists = errors.New("file already exists")

func Send(filename string) error {
	tpl, err := template.Load(filename)
	if err != nil {
		return err
	}

	req := RequestFromTemplate(tpl)

	client := httpclient.New()

	resp, err := client.Send(req)
	if err != nil {
		return err
	}

	output.Print(resp)

	return nil
}

func RequestFromTemplate(tpl *template.Template) httpclient.Request {
	return httpclient.Request{
		Method:  string(tpl.Request.Method),
		Url:     tpl.Request.Url,
		Headers: tpl.Request.Headers,
		Body:    tpl.Request.Body.Template,
	}
}
