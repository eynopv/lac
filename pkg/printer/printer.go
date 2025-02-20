package printer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PrinterConfig struct {
	PrintResponseBody    bool
	PrintResponseHeaders bool
	PrintRequestBody     bool
	PrintRequestHeaders  bool
}

type Printer struct {
	config PrinterConfig
}

func NewPrinter(config PrinterConfig) Printer {
	return Printer{
		config: config,
	}
}

func PrintHeaders(headers http.Header) {
	fmt.Print(StringifyHeaders(headers))
	fmt.Print("\n")
}

func StringifyHeaders(headers http.Header) string {
	var s []string
	for key, value := range headers {
		s = append(s, fmt.Sprintf("%s: %s", Cyan(key), strings.Join(value, ", ")))
	}

	return strings.Join(s, "\n")
}

func PrintPrettyJson(v any) {
	prettyJson, err := ToPrettyJsonString(v)
	if err != nil {
		fmt.Printf("Failed to parse %v to json: %v\n", v, err)
		return
	}

	fmt.Println(prettyJson)
}

func ToPrettyJsonString(v any) (string, error) {
	prettyJson, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return "", err
	}

	return string(prettyJson), nil
}

func Red(s string) string {
	return ColorRed + s + ColorReset
}

func Green(s string) string {
	return ColorGreen + s + ColorReset
}

func Yellow(s string) string {
	return ColorYellow + s + ColorReset
}

func Blue(s string) string {
	return ColorBlue + s + ColorReset
}

func Magenta(s string) string {
	return ColorMagenta + s + ColorReset
}

func Cyan(s string) string {
	return ColorCyan + s + ColorReset
}
