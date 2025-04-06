package utils

import (
	"testing"

	"github.com/eynopv/lac/internal/assert"
)

func TestFlattenMap(t *testing.T) {
	m := map[string]any{
		"a": map[string]any{
			"b": "c",
			"d": map[string]any{
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

	assert.DeepEqual(t, result, expected)

	expected = map[string]string{
		"prefix.a.b":   "c",
		"prefix.a.d.e": "f",
		"prefix.g":     "k",
	}
	result = FlattenMap(m, "prefix")
	assert.DeepEqual(t, result, expected)
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

	assert.Equal(t, m["shouldNotOverwrite"], m1["shouldNotOverwrite"])
	assert.Equal(t, m["shouldOverwrite"], m2["shouldOverwrite"])
}
