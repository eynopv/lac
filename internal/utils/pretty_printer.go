package utils

import (
	"encoding/json"
	"fmt"
)

func PrintPrettyJson(v any) {
	prettyJson, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Printf("Failed to parse %v to json: %v\n", v, err)
		return
	}
	fmt.Println(string(prettyJson))
}
