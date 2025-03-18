package validators

import "encoding/json"

func IsJson(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}
