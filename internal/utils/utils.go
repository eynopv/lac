package utils

import "fmt"

func FlattenMap(input map[string]interface{}, prefix string) map[string]string {
	flattened := map[string]string{}

	for key, value := range input {
		var newKey string
		if prefix == "" {
			newKey = key
		} else {
			newKey = fmt.Sprintf("%s.%s", prefix, key)
		}

		switch child := value.(type) {
		case map[string]interface{}:
			submap := FlattenMap(child, newKey)
			for k, v := range submap {
				flattened[k] = v
			}
		default:
			flattened[newKey] = fmt.Sprintf("%v", value)
		}
	}

	return flattened
}
