package printer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eynopv/lac/pkg/result"
)

type Formatter interface {
	Format(res *result.Result) (string, error)
}

type JsonFormatter struct {
	includes Includes
}

func (f *JsonFormatter) Format(res *result.Result) (string, error) {
	m := map[string]any{}

	if f.includes.RequestHeaders || f.includes.RequestBody || f.includes.RequestMeta {
		rm := map[string]any{}

		if f.includes.RequestMeta {
			rm["meta"] = res.RequestLine()
		}

		if f.includes.RequestHeaders {
			rm["headers"] = res.Response.Request.Header
		}

		if f.includes.RequestBody {
			rm["body"] = res.RequestBody.Data()
		}

		m["request"] = rm
	}

	if f.includes.ResponseHeaders || f.includes.ResponseBody || f.includes.ResponseMeta {
		rm := map[string]any{}

		if f.includes.ResponseMeta {
			rm["meta"] = res.StatusLine()
		}

		if f.includes.ResponseHeaders {
			rm["headers"] = res.Response.Header
		}

		if f.includes.ResponseBody {
			rm["body"] = res.ResponseBody.Data()
		}

		m["response"] = rm
	}

	b, err := json.MarshalIndent(m, "", "  ")

	if err != nil {
		return "", err
	}

	return string(b), nil
}

type PrettyFormatter struct {
	includes Includes
}

func (f *PrettyFormatter) Format(res *result.Result) (string, error) {
	sections := []string{}

	if f.includes.RequestHeaders {
		sections = append(sections, f.printRequestHeaders(res))
	}

	if f.includes.RequestBody {
		sections = append(sections, f.printBody(&res.RequestBody))
	}

	if f.includes.ResponseHeaders {
		sections = append(sections, f.printResponseHeaders(res))
	}

	if f.includes.ResponseBody {
		sections = append(sections, f.printBody(&res.ResponseBody))
	}

	return strings.Join(sections, "\n"), nil
}

func (f *PrettyFormatter) printRequestHeaders(res *result.Result) string {
	req := *res.Response.Request

	if f.includes.RequestMeta {
		return f.requestLine(*res.RequestLine()) + f.headers(req.Header)
	}

	return f.headers(req.Header)
}

func (f *PrettyFormatter) printResponseHeaders(res *result.Result) string {
	if f.includes.ResponseMeta {
		return f.statusLine(*res.StatusLine()) + f.headers(res.Response.Header)
	}

	return f.headers(res.Response.Header)
}

func (f *PrettyFormatter) printBody(body *result.Body) string {
	if jsonBody := body.Json(); jsonBody != nil {
		return fmt.Sprintf("%v\n", f.json(jsonBody))
	}

	if textBody := body.Text(); textBody != "" {
		return fmt.Sprintf("%v\n", textBody)
	}

	return ""
}

func (f *PrettyFormatter) headers(headers http.Header) string {
	return formatHeaders(headers, true)
}

func (f *PrettyFormatter) statusLine(line result.StatusLine) string {
	return formatStatusLine(line, true)
}

func (f *PrettyFormatter) requestLine(line result.RequestLine) string {
	return formatRequestLine(line, true)
}

func (f *PrettyFormatter) json(j map[string]any) string {
	return formatJson(j, true)
}

func formatHeaders(headers http.Header, colorized bool) string {
	fstring := "%s: %s\n"

	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var sb strings.Builder

	for _, key := range keys {
		value := strings.Join(headers[key], ", ")
		if colorized {
			sb.WriteString(fmt.Sprintf(fstring, Cyan(key), value))
		} else {
			sb.WriteString(fmt.Sprintf(fstring, key, value))
		}
	}

	return sb.String()
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

func formatJson(j map[string]any, colorized bool) string {
	var f func(inner map[string]any, level int) string

	var formatValue func(value any, level int) string

	f = func(inner map[string]any, level int) string {
		sb := strings.Builder{}
		sb.WriteString("{")

		keys := []string{}
		for k := range inner {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		padding := strings.Repeat("  ", level+1)

		for i, k := range keys {
			sb.WriteString("\n")

			value := inner[k]
			isLast := i == len(keys)-1

			sb.WriteString(padding)
			sb.WriteString(formatJsonKey(k, colorized))

			sb.WriteString(formatValue(value, level))

			if !isLast {
				sb.WriteString(",")
			} else {
				sb.WriteString("\n")
			}
		}

		padding = strings.Repeat("  ", level)
		if sb.Len() != 1 {
			sb.WriteString(padding)
		}

		sb.WriteString("}")

		return sb.String()
	}

	formatValue = func(value any, level int) string {
		switch v := value.(type) {
		case string:
			return formatJsonStringValue(v, colorized)
		case int, float64:
			return formatJsonNumberValue(v, colorized)
		case bool:
			return formatJsonBoolValue(v, colorized)
		case nil:
			return formatJsonNilValue(colorized)
		case map[string]any:
			return f(v, level+1)
		case []any:
			sb := strings.Builder{}
			sb.WriteString("[")

			for i, vv := range v {
				sb.WriteString(formatValue(vv, level))

				if i != len(v)-1 {
					sb.WriteString(", ")
				}
			}

			sb.WriteString("]")

			return sb.String()
		}

		return ""
	}

	return f(j, 0)
}

func formatJsonKey(s string, colorized bool) string {
	fstring := `%s: `
	value := strconv.Quote(s)

	if !colorized {
		return fmt.Sprintf(fstring, value)
	}

	return fmt.Sprintf(fstring, Yellow(value))
}

func formatJsonStringValue(s string, colorized bool) string {
	if !colorized {
		return strconv.Quote(s)
	}

	return Green(strconv.Quote(s))
}

func formatJsonNumberValue(n any, colorized bool) string {
	s := fmt.Sprintf("%v", n)
	if !colorized {
		return s
	}

	return Cyan(s)
}

func formatJsonBoolValue(b bool, colorized bool) string {
	s := fmt.Sprintf("%v", b)
	if !colorized {
		return s
	}

	return Magenta(s)
}

func formatJsonNilValue(colorized bool) string {
	if !colorized {
		return "null"
	}

	return Red("null")
}
