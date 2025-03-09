package includes

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func TestRemoteInclude_Type(t *testing.T) {
	remoteInclude := &RemoteInclude{}
	expectedType := "remote"
	if remoteInclude.Type() != expectedType {
		t.Fatalf("Type() = %v, want %v", remoteInclude.Type(), expectedType)
	}
}

func TestRemoteInclude_Equals(t *testing.T) {
	tests := []struct {
		name     string
		include1 *RemoteInclude
		include2 Include
		expected bool
	}{
		{
			name: "equal includes",
			include1: &RemoteInclude{
				Remote: "https://example.com/test.yml",
			},
			include2: &RemoteInclude{
				Remote: "https://example.com/test.yml",
			},
			expected: true,
		},
		{
			name: "different remotes",
			include1: &RemoteInclude{
				Remote: "https://example.com/test1.yml",
			},
			include2: &RemoteInclude{
				Remote: "https://example.com/test2.yml",
			},
			expected: false,
		},
		{
			name: "different types",
			include1: &RemoteInclude{
				Remote: "https://example.com/test.yml",
			},
			include2: &LocalInclude{
				Path: "test.yml",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.include1.Equals(tt.include2)
			if result != tt.expected {
				t.Fatalf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewRemoteInclude(t *testing.T) {
	node := &yaml.Node{}
	remote := "https://example.com/test.yml"
	integrity := "sha256-abcdef"
	remoteInclude := NewRemoteInclude(node, remote, integrity)

	if remoteInclude.Remote != remote {
		t.Fatalf("NewRemoteInclude().Remote = %v, want %v", remoteInclude.Remote, remote)
	}

	if remoteInclude.Integrity != integrity {
		t.Fatalf("NewRemoteInclude().Integrity = %v, want %v", remoteInclude.Integrity, integrity)
	}

	if remoteInclude.Node() != node {
		t.Fatalf("NewRemoteInclude().Node() = %v, want %v", remoteInclude.Node(), node)
	}
}
