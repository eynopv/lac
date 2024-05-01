package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFlattenMap(t *testing.T) {
	m := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "c",
			"d": map[string]interface{}{
				"e": "f",
			},
		},
		"g": "k",
	}

	expected := map[string]string{
		"a.b":   "c",
		"a.d.e": "f",
		"g":     "k",
	}
	result := FlattenMap(m, "")
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("\nExpected: " + fmt.Sprint(expected) + "\nReceived: " + fmt.Sprint(result) + "\n")
	}

	expected = map[string]string{
		"prefix.a.b":   "c",
		"prefix.a.d.e": "f",
		"prefix.g":     "k",
	}
	result = FlattenMap(m, "prefix")
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("\nExpected: " + fmt.Sprint(expected) + "\nReceived: " + fmt.Sprint(result) + "\n")
	}
}
