package utils

import (
	"fmt"
	"strconv"
)

func ToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case float64:
		// Check if it's an integer
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10), nil
		} else {
			return strconv.FormatFloat(v, 'f', -1, 64), nil
		}
	case nil:
		return "null", nil
	}

	return "", fmt.Errorf("unable to convert to string")
}
