package utils

import (
	"net/http"
	"strings"
)

func StringToHttpMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return http.MethodGet
	case "POST":
		return http.MethodPost
	case "PUT":
		return http.MethodPut
	case "PATCH":
		return http.MethodPatch
	case "DELETE":
		return http.MethodDelete
	default:
		return "UNKNOWN METHOD"
	}
}
