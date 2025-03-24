package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

var ErrBasicAuthParse = errors.New("failed to parse basic auth")
var ErrBasicAuthInvalid = errors.New("username and password are required")

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
		return nil, fmt.Errorf("%w: %v", ErrBasicAuthParse, err)
	}

	if wrapper.Auth != nil && (wrapper.Auth.Username == "" || wrapper.Auth.Password == "") {
		return nil, ErrBasicAuthInvalid
	}

	return wrapper.Auth, nil
}

func (bs *BasicAuth) Apply(r *http.Request) {
	r.SetBasicAuth(bs.Username, bs.Password)
}
