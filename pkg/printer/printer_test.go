package printer

import (
	"testing"
)

func TestToPrettyJsonString(t *testing.T) {
	m := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "c",
			"d": map[string]interface{}{
				"e": "f",
			},
		},
		"g": "k",
	}

	_, err := ToPrettyJsonString(m)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
