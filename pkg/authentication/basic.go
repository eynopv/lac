package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v3"

	"github.com/eynopv/lac/internal/errorsx"
	"github.com/eynopv/lac/pkg/template"
)

type BasicAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func NewBasicAuth(t *template.Template) (*BasicAuth, error) {
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
