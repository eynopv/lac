package printer

import (
	"encoding/json"
	"fmt"
)

func PrintPrettyJson(v any) {
	prettyJson, err := toPrettyJsonString(v)
	if err != nil {
		fmt.Printf("Failed to parse %v to json: %v\n", v, err)
		return
	}
	fmt.Println(prettyJson)
}

func toPrettyJsonString(v any) (string, error) {
	prettyJson, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return "", err
	}
	return string(prettyJson), nil
}

func Red(s string) string {
	return colorRed + s + colorReset
}

func Green(s string) string {
	return colorGreen + s + colorReset
}

func Yellow(s string) string {
	return colorYellow + s + colorReset
}

func Blue(s string) string {
	return colorBlue + s + colorReset
}

func Magenta(s string) string {
	return colorMagenta + s + colorReset
}

func Cyan(s string) string {
	return colorCyan + s + colorReset
}
