package request

import (
	"encoding/json"
	"errors"
	"fmt"

	yaml "gopkg.in/yaml.v3"

	"github.com/eynopv/lac/pkg/param"
	"github.com/eynopv/lac/pkg/utils"
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

func (t *Template) Interpolate(vars map[string]any, useEnv bool) *Template {
	result := Template(param.Param(*t).Resolve(vars, true))
	return &result
}

func (t *Template) Parse() (*Request, error) {
	var requestData RequestData

	err := json.Unmarshal([]byte(*t), &requestData)
	if err != nil {
		err = yaml.Unmarshal([]byte(*t), &requestData)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTemplateParse, err)
	}

	nr := NewRequest(requestData)

	return &nr, nil
}
