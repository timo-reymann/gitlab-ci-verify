package format_conversion

import (
	"testing"
)

func TestParseYamlNode(t *testing.T) {
	_, err := ParseYamlNode([]byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
}
