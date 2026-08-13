package app

import (
	"github.com/eynopv/lac/internal/template"
)

func New(filename string) error {
	t := template.Template{}

	template.Save(filename, t)

	return nil
}
