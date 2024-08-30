package format_conversion

import "encoding/json"

// ParseJson from the given raw content into a generic map
func ParseJson(content []byte) (map[string]any, error) {
	result := map[string]any{}
	err := json.Unmarshal(content, &result)
	return result, err
}

func ToJson(content map[string]any) ([]byte, error) {
	return json.Marshal(content)
}
