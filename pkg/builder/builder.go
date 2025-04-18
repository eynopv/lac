package builder

import (
	"github.com/eynopv/lac/pkg/authentication"
	"github.com/eynopv/lac/pkg/client"
	"github.com/eynopv/lac/pkg/printer"
	"github.com/eynopv/lac/pkg/request"
	"github.com/eynopv/lac/pkg/template"
	"github.com/eynopv/lac/pkg/utils"
	"github.com/eynopv/lac/pkg/variables"
)

type Builder struct {
	ClientConfig client.ClientConfig
	TemplatePath string
	Variables    variables.Variables
	Headers      map[string]request.StringOrStringList
}

func (b *Builder) BuildClient() *client.Client {
	return client.NewClient(&b.ClientConfig)
}

func (b *Builder) BuildTemplate() (*template.Template, error) {
	return template.NewTemplate(b.TemplatePath)
}

func (b *Builder) BuildRequest() (*request.Request, error) {
	t, err := b.BuildTemplate()
	if err != nil {
		return nil, err
	}

	t = t.Interpolate(b.Variables, true)

	r, err := t.Parse()
	if err != nil {
		return nil, err
	}

	r.Headers = utils.CombineMaps(b.Headers, r.Headers)

	return r, nil
}

func (b *Builder) BuildAuth() (authentication.Auth, error) {
	t, err := b.BuildTemplate()
	if err != nil {
		return nil, err
	}

	t = t.Interpolate(b.Variables, true)

	return authentication.NewAuth(t)
}

func (b *Builder) BuildPrinter() *printer.Printer {
	p := printer.NewPrinter(b.ClientConfig.PrinterConfig)
	return &p
}
