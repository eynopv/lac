package request

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

type BasicAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func NewBasicAuth(t *Template) *BasicAuth {
	var basicAuth BasicAuth

	bs := []byte(*t)

	if err := json.Unmarshal(bs, &basicAuth); err == nil {
		return &basicAuth
	}

	if err := yaml.Unmarshal(bs, &basicAuth); err == nil {
		return &basicAuth
	}

	return nil
}

func (bs *BasicAuth) Apply(r *http.Request) {
	r.SetBasicAuth(bs.Username, bs.Password)
}
