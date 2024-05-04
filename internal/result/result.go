package result

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eynopv/lac/internal/utils/printer"
)

type Result struct {
	Status      string
	StatusCode  int
	ElapsedTime time.Duration
	Body        map[string]interface{}
	Text        string
	Headers     http.Header
}

func NewResult(
	elapsedTime time.Duration,
	status string,
	statusCode int,
	headers http.Header,
	bodyRaw []byte,
) (Result, error) {
	result := Result{
		ElapsedTime: elapsedTime,
		Status:      status,
		StatusCode:  statusCode,
		Headers:     headers,
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

func (r *Result) Print(showHeaders bool) {
	if r.StatusCode < 300 {
		fmt.Println(printer.Green(r.Status))
	} else {
		fmt.Println(printer.Red(r.Status))
	}

	fmt.Println(fmt.Sprintf("%s: %s", printer.Cyan("Elapsed Time"), r.ElapsedTime))

	if showHeaders {
		for key, value := range r.Headers {
			fmt.Println(fmt.Sprintf("%s: %s", printer.Cyan(key), strings.Join(value, ", ")))
		}
	}

	if r.Body != nil {
		fmt.Println()
		printer.PrintPrettyJson(r.Body)
	} else if r.Text != "" {
		fmt.Println()
		fmt.Println(r.Text)
	}
}
