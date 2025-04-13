package printer

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eynopv/lac/pkg/result"
)

type Formatter struct {
	colored bool
}

func (f Formatter) Headers(headers http.Header) string {
	return formatHeaders(headers, f.colored)
}

func (f Formatter) StatusLine(line result.StatusLine) string {
	return formatStatusLine(line, f.colored)
}

func (f Formatter) RequestLine(line result.RequestLine) string {
	return formatRequestLine(line, f.colored)
}

func (f Formatter) Json(j map[string]any) string {
	return formatJson(j, f.colored)
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
