package http_method

import (
	"net/http"
	"strings"
)

var methodMap = map[string]string{
	"GET":    http.MethodGet,
	"POST":   http.MethodPost,
	"PUT":    http.MethodPut,
	"PATCH":  http.MethodPatch,
	"DELETE": http.MethodDelete,
}

func NormalizeHttpMethod(method string) string {
	uppercaseMethod := strings.ToUpper(method)

	if value, ok := methodMap[uppercaseMethod]; ok {
		return value
	}

	return "UNKNOWN"
}
