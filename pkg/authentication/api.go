package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v3"

	"github.com/eynopv/lac/internal/errorsx"
	"github.com/eynopv/lac/pkg/template"
)

type ApiAuth struct {
	Header string `json:"header" yaml:"header"`
	Key    string `json:"key" yaml:"key"`
}

func NewApiAuth(t *template.Template) (*ApiAuth, error) {
	var wrapper struct {
		Auth *ApiAuth `json:"auth" yaml:"auth"`
	}

	bs := []byte(*t)

	err := json.Unmarshal(bs, &wrapper)
	if err != nil {
		err = yaml.Unmarshal(bs, &wrapper)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errorsx.ErrApiAuthParse, err)
	}

	if wrapper.Auth != nil && (wrapper.Auth.Header == "" || wrapper.Auth.Key == "") {
		return nil, errorsx.ErrApiAuthInvalid
	}

	return wrapper.Auth, nil
}

func (a *ApiAuth) Apply(r *http.Request) {
	r.Header.Set(a.Header, a.Key)
}

func (a *ApiAuth) GetType() AuthType {
	return Api
}
