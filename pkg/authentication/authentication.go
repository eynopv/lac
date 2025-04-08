package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v3"

	"github.com/eynopv/lac/internal/errorsx"
	"github.com/eynopv/lac/pkg/template"
)

type AuthType string

const (
	Api    AuthType = "api"
	Basic  AuthType = "basic"
	Bearer AuthType = "bearer"
)

type Auth interface {
	Apply(r *http.Request)
	GetType() AuthType
}

type AuthBase struct {
	Type AuthType `json:"type" yaml:"type"`
}

func NewAuth(t *template.Template) (Auth, error) {
	var wrapper struct {
		Auth *AuthBase `json:"auth" yaml:"auth"`
	}

	bs := []byte(*t)

	err := json.Unmarshal(bs, &wrapper)
	if err != nil {
		if err = yaml.Unmarshal(bs, &wrapper); err != nil {
			return nil, fmt.Errorf("%w: %v", errorsx.ErrAuthParse, err)
		}
	}

	if wrapper.Auth == nil {
		return nil, nil
	}

	switch wrapper.Auth.Type {
	case Api:
		return NewApiAuth(t)
	case Basic:
		return NewBasicAuth(t)
	case Bearer:
		return NewBearerAuth(t)
	default:
		return nil, errorsx.ErrAuthUnknown
	}
}
