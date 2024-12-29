package hashing

import (
	"testing"
)

func TestCreateHashFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello world", "hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
		{"Go programming", "Go programming", "083dcfbba2677b06ecf81951f97fdb73eca905591285f30ca07911bced99dbe4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateHashFromString(tt.input)
			if result != tt.expected {
				t.Errorf("CreateHashFromString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
