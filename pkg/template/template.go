package template

import (
	"encoding/json"
	"errors"
	"fmt"

	yaml "gopkg.in/yaml.v3"

	"github.com/eynopv/lac/pkg/param"
	"github.com/eynopv/lac/pkg/request"
	"github.com/eynopv/lac/pkg/utils"
	"github.com/eynopv/lac/pkg/variables"
)

var ErrTemplateParse = errors.New("failed to parse template")

type Template string

var fileLoader = utils.LoadFile

func NewTemplate(templatePath string) (*Template, error) {
	data, err := fileLoader(templatePath)
	if err != nil {
		return nil, err
	}

	result := Template(*data)

	return &result, nil
}

func (t *Template) Interpolate(vars variables.Variables, useEnv bool) *Template {
	result := Template(param.Param(*t).Resolve(vars, true))
	return &result
}

func (t *Template) Parse() (*request.Request, error) {
	var requestData request.RequestData

	err := json.Unmarshal([]byte(*t), &requestData)
	if err != nil {
		err = yaml.Unmarshal([]byte(*t), &requestData)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTemplateParse, err)
	}

	nr := request.NewRequest(requestData)

	return &nr, nil
}
