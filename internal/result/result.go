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
	fmt.Println("Status:", r.Status)
	fmt.Println("Elapsed Time:", r.ElapsedTime)

	if showHeaders {
		printer.PrintPrettyJson(r.Headers)
	}

	if r.Body != nil {
		printer.PrintPrettyJson(r.Body)
	} else if r.Text != "" {
		fmt.Println(r.Text)
	}
}
