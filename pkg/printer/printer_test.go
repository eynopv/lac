package printer

import (
	"testing"

	"github.com/eynopv/lac/internal/assert"
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
	assert.NoError(t, err)
}
