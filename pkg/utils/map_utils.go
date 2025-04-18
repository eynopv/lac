package utils

import (
	"fmt"
	"maps"
)

func FlattenMap(input map[string]any, prefix string) map[string]string {
	flattened := map[string]string{}

	for key, value := range input {
		var newKey string
		if prefix == "" {
			newKey = key
		} else {
			newKey = fmt.Sprintf("%s.%s", prefix, key)
		}

		switch child := value.(type) {
		case map[string]any:
			submap := FlattenMap(child, newKey)
			maps.Copy(flattened, submap)
		default:
			flattened[newKey] = fmt.Sprintf("%v", value)
		}
	}

	return flattened
}

func CombineMaps[T any](input ...map[string]T) map[string]T {
	finalMap := map[string]T{}
	for _, m := range input {
		maps.Copy(finalMap, m)
	}

	return finalMap
}
