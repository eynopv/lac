package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"

	"github.com/eynopv/lac/internal/errorsx"
)

type AuthType string

const (
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

func NewAuth(t *Template) (Auth, error) {
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
	case Basic:
		return NewBasicAuth(t)
	case Bearer:
		return NewBearerAuth(t)
	default:
		return nil, errorsx.ErrAuthUnknown
	}
}

type BasicAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func NewBasicAuth(t *Template) (*BasicAuth, error) {
	var wrapper struct {
		Auth *BasicAuth `json:"auth" yaml:"auth"`
	}

	bs := []byte(*t)

	err := json.Unmarshal(bs, &wrapper)
	if err != nil {
		err = yaml.Unmarshal(bs, &wrapper)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errorsx.ErrBasicAuthParse, err)
	}

	if wrapper.Auth != nil && (wrapper.Auth.Username == "" || wrapper.Auth.Password == "") {
		return nil, errorsx.ErrBasicAuthInvalid
	}

	return wrapper.Auth, nil
}

func (a *BasicAuth) Apply(r *http.Request) {
	if a != nil {
		r.SetBasicAuth(a.Username, a.Password)
	}
}

func (a *BasicAuth) GetType() AuthType {
	return Basic
}

type BearerAuth struct {
	Token string `json:"token"`
}

func NewBearerAuth(t *Template) (*BearerAuth, error) {
	var wrapper struct {
		Auth *BearerAuth `json:"auth" yaml:"auth"`
	}

	bs := []byte(*t)

	err := json.Unmarshal(bs, &wrapper)
	if err != nil {
		if err = yaml.Unmarshal(bs, &wrapper); err != nil {
			return nil, fmt.Errorf("%w: %v", errorsx.ErrBearerAuthParse, err)
		}
	}

	if wrapper.Auth != nil && wrapper.Auth.Token == "" {
		return nil, errorsx.ErrBearerAuthInvalid
	}

	return wrapper.Auth, nil
}

func (a *BearerAuth) GetType() AuthType {
	return Bearer
}

func (a *BearerAuth) Apply(r *http.Request) {
	if a != nil {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", a.Token))
	}
}
