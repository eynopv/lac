package printer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/eynopv/lac/pkg/result"
)

type Formatter interface {
	StatusLine(statusLine result.StatusLine) string
	RequestLine(requestLine result.RequestLine) string
	Headers(headers http.Header) string
	Json(map[string]any) string
}

func NewFormatter(isTerminal bool) Formatter {
	if isTerminal {
		return ColorFormatter{}
	}

	return PlainFormatter{}
}

type ColorFormatter struct{}
type PlainFormatter struct{}

func (f ColorFormatter) Headers(headers http.Header) string {
	return formatHeaders(headers, true)
}

func (f PlainFormatter) Headers(headers http.Header) string {
	return formatHeaders(headers, false)
}

func (f ColorFormatter) StatusLine(line result.StatusLine) string {
	return formatStatusLine(line, true)
}

func (f PlainFormatter) StatusLine(line result.StatusLine) string {
	return formatStatusLine(line, false)
}

func (f ColorFormatter) RequestLine(line result.RequestLine) string {
	return formatRequestLine(line, true)
}

func (f PlainFormatter) RequestLine(line result.RequestLine) string {
	return formatRequestLine(line, false)
}

func (f ColorFormatter) Json(j map[string]interface{}) string {
	return formatJson(j)
}

func (f PlainFormatter) Json(j map[string]interface{}) string {
	return formatJson(j)
}

func formatHeaders(headers http.Header, colorized bool) string {
	s := ""

	fstring := "%s: %s\n"

	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		value := strings.Join(headers[key], ", ")
		if colorized {
			s += fmt.Sprintf(fstring, Cyan(key), value)
		} else {
			s += fmt.Sprintf(fstring, key, value)
		}
	}

	return s
}

func formatStatusLine(line result.StatusLine, colorized bool) string {
	fstring := "%v %v [%v]\n"

	if !colorized {
		return fmt.Sprintf(fstring, line.Protocol, line.Status, line.Time)
	}

	var timeColor, statusColor Color

	switch {
	case strings.HasPrefix(line.Status, "2"), strings.HasPrefix(line.Status, "1"):
		statusColor = ColorGreen
	case strings.HasPrefix(line.Status, "3"):
		statusColor = ColorCyan
	default:
		statusColor = ColorRed
	}

	switch {
	case line.Time < 500*time.Millisecond:
		timeColor = ColorReset
	case line.Time < 1000*time.Millisecond:
		timeColor = ColorYellow
	default:
		timeColor = ColorRed
	}

	return fmt.Sprintf(
		fstring,
		line.Protocol,
		Colorize(line.Status, statusColor),
		Colorize(line.Time.String(), timeColor),
	)
}

func formatRequestLine(line result.RequestLine, colorized bool) string {
	fstring := "%v %v %v\n"

	if !colorized {
		return fmt.Sprintf(fstring, line.Method, line.Url, line.Protocol)
	}

	var methodColor Color

	switch line.Method {
	case http.MethodGet, http.MethodHead:
		methodColor = ColorGreen
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		methodColor = ColorYellow
	case http.MethodDelete:
		methodColor = ColorRed
	default:
		methodColor = ColorMagenta
	}

	return fmt.Sprintf(fstring, Colorize(line.Method, methodColor), line.Url, line.Protocol)
}

func formatJson(j map[string]any) string {
	if prettyJson, err := json.MarshalIndent(j, "", " "); err != nil {
		return "unable to parse json"
	} else {
		return string(prettyJson)
	}
}
