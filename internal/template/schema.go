package template

import (
	"net/http"
)

type Template struct {
	Request Request `json:"request"`
}

type Request struct {
	Method    Method            `json:"method"`
	Url       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      Body              `json:"body"`
	Variables map[string]any    `json:"variables"`
}

type Method string

const (
	MethodGet     Method = http.MethodGet
	MethodHead    Method = http.MethodHead
	MethodPost    Method = http.MethodPost
	MethodPut     Method = http.MethodPut
	MethodPatch   Method = http.MethodPatch
	MethodDelete  Method = http.MethodDelete
	MethodConnect Method = http.MethodConnect
	MethodOptions Method = http.MethodOptions
	MethodTrace   Method = http.MethodTrace
)

type Body struct {
	Template string `json:"template"`
}
