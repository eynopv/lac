package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"

	"github.com/eynopv/lac/internal/errorsx"
	"github.com/eynopv/lac/pkg/request"
)

type BearerAuth struct {
	Token string `json:"token" yaml:"token"`
}

func NewBearerAuth(t *request.Template) (*BearerAuth, error) {
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
