package result

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eynopv/lac/pkg/printer"
)

type Result struct {
	Path        string
	Status      string
	StatusCode  int
	Protocol    string
	ElapsedTime time.Duration
	Body        map[string]interface{}
	Text        string
	Headers     http.Header
}

func NewResult(
	elapsedTime time.Duration,
	path string,
	status string,
	statusCode int,
	headers http.Header,
	protocol string,
	bodyRaw []byte,
) (Result, error) {
	result := Result{
		Path:        path,
		ElapsedTime: elapsedTime,
		Status:      status,
		StatusCode:  statusCode,
		Headers:     headers,
		Protocol:    protocol,
	}

	if len(bodyRaw) > 0 {
		contentType := headers.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			var responseData map[string]interface{}
			err := json.Unmarshal(bodyRaw, &responseData)
			if err != nil {
				return result, err
			}
			result.Body = responseData
		} else if strings.Contains(contentType, "text/") {
			result.Text = string(bodyRaw)
		}
	}

	return result, nil
}

func (r *Result) Print() {
	fmt.Println(r.Path)

	if r.StatusCode < 300 {
		fmt.Printf("%v %v [%v]\n", r.Protocol, printer.Green(r.Status), r.ElapsedTime)
	} else if r.StatusCode >= 300 && r.StatusCode < 400 {
		fmt.Printf("%v %v [%v]\n", r.Protocol, printer.Cyan(r.Status), r.ElapsedTime)
	} else {
		fmt.Printf("%v %v [%v]\n", r.Protocol, printer.Red(r.Status), r.ElapsedTime)
	}

	fmt.Println()

	for key, value := range r.Headers {
		fmt.Printf("%s: %s\n", printer.Cyan(key), strings.Join(value, ", "))
	}

	if r.Body != nil {
		fmt.Println()
		printer.PrintPrettyJson(r.Body)
	} else if r.Text != "" {
		fmt.Println()
		fmt.Println(r.Text)
	}
}
