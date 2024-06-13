package result

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eynopv/lac/pkg/utils/printer"
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
	fmt.Println()

	if r.StatusCode < 300 {
		fmt.Println(r.Protocol, printer.Green(r.Status))
	} else if r.StatusCode >= 300 && r.StatusCode < 400 {
		fmt.Println(r.Protocol, printer.Cyan(r.Status))
	} else {
		fmt.Println(r.Protocol, printer.Red(r.Status))
	}

	fmt.Println(fmt.Sprintf("%s: %s", printer.Cyan("Elapsed Time"), r.ElapsedTime))

	for key, value := range r.Headers {
		fmt.Println(fmt.Sprintf("%s: %s", printer.Cyan(key), strings.Join(value, ", ")))
	}

	if r.Body != nil {
		fmt.Println()
		printer.PrintPrettyJson(r.Body)
	} else if r.Text != "" {
		fmt.Println()
		fmt.Println(r.Text)
	}
}
