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

func TestCombineMaps(t *testing.T) {
	m1 := map[string]string{
		"shouldNotOverwrite": "shouldNotOverwrite",
		"shouldOverwrite":    "notOverwritten",
	}
	m2 := map[string]string{
		"shouldOverwrite": "overwritten",
	}
	m := CombineMaps(m1, m2)
	if m["shouldNotOverwrite"] != "shouldNotOverwrite" {
		t.Fatalf("Expected field to no be overwritten")
	}
	if m["shouldOverwrite"] == m1["shouldOverwrite"] {
		t.Fatalf("Expected field to be overwritten")
	}
}
