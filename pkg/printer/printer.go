package printer

import (
	"encoding/json"
	"fmt"
)

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
