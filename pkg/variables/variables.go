package variables

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v3"

	"github.com/eynopv/lac/internal/errorsx"
)

type Variables map[string]any

func (v *Variables) UnmarshalJSON(data []byte) error {
	var temp map[string]any
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	for key, value := range temp {
		switch value.(type) {
		case string, float64, bool, nil:
			// allowed
		default:
			return fmt.Errorf("%w: key %s", errorsx.ErrUnsupportedVariablesValue, key)
		}
	}

	*v = temp

	return nil
}

func (v *Variables) UnmarshalYAML(value *yaml.Node) error {
	var raw any
	if err := value.Decode(&raw); err != nil {
		return err
	}

	jsonData, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	return v.UnmarshalJSON(jsonData)
}
