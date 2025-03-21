package includes

import (
	"testing"
)

func TestFilterByTypes(t *testing.T) {
	tests := []struct {
		name         string
		includes     []Include
		includeTypes []string
		expected     []Include
	}{
		{
			name: "filter by single type",
			includes: []Include{
				&LocalInclude{Path: "test1.yml"},
				&ProjectInclude{Project: "my-group/my-project", Files: []string{"test2.yml"}},
			},
			includeTypes: []string{"local"},
			expected: []Include{
				&LocalInclude{Path: "test1.yml"},
			},
		},
		{
			name: "filter by multiple types",
			includes: []Include{
				&LocalInclude{Path: "test1.yml"},
				&ProjectInclude{Project: "my-group/my-project", Files: []string{"test2.yml"}},
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
			includeTypes: []string{"local", "remote"},
			expected: []Include{
				&LocalInclude{Path: "test1.yml"},
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
		},
		{
			name: "no matching types",
			includes: []Include{
				&ProjectInclude{Project: "my-group/my-project", Files: []string{"test2.yml"}},
			},
			includeTypes: []string{"local"},
			expected:     []Include{},
		},
		{
			name:         "empty includes",
			includes:     []Include{},
			includeTypes: []string{"local"},
			expected:     []Include{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterByTypes(tt.includes, tt.includeTypes...)
			if !Equals(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		name     string
		a        []Include
		b        []Include
		expected bool
	}{
		{
			name: "equal slices",
			a: []Include{
				&LocalInclude{Path: "test1.yml"},
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
			b: []Include{
				&LocalInclude{Path: "test1.yml"},
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
			expected: true,
		},
		{
			name: "different slices",
			a: []Include{
				&LocalInclude{Path: "test1.yml"},
			},
			b: []Include{
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
			expected: false,
		},
		{
			name: "different lengths",
			a: []Include{
				&LocalInclude{Path: "test1.yml"},
			},
			b: []Include{
				&LocalInclude{Path: "test1.yml"},
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equals(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		includes []Include
		expected []Include
	}{
		{
			name: "unique elements",
			includes: []Include{
				&LocalInclude{Path: "test1.yml"},
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
			expected: []Include{
				&LocalInclude{Path: "test1.yml"},
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
		},
		{
			name: "duplicate elements",
			includes: []Include{
				&LocalInclude{Path: "test1.yml"},
				&LocalInclude{Path: "test1.yml"},
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
			expected: []Include{
				&LocalInclude{Path: "test1.yml"},
				&RemoteInclude{Remote: "https://example.com/test.yml"},
			},
		},
		{
			name: "all duplicates",
			includes: []Include{
				&LocalInclude{Path: "test1.yml"},
				&LocalInclude{Path: "test1.yml"},
				&LocalInclude{Path: "test1.yml"},
			},
			expected: []Include{
				&LocalInclude{Path: "test1.yml"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.includes)
			if !Equals(result, tt.expected) {
				t.Errorf("Unique() = %v, want %v", result, tt.expected)
			}
		})
	}
}
