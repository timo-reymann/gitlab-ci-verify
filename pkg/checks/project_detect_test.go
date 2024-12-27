package checks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHasProjectPoliciesOnDisk(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "No .gitlab-ci-verify directory",
			setup: func(dir string) error {
				return nil
			},
			expected: false,
		},
		{
			name: "No checks directory",
			setup: func(dir string) error {
				return os.Mkdir(filepath.Join(dir, ".gitlab-ci-verify"), 0755)
			},
			expected: false,
		},
		{
			name: "Checks directory exists",
			setup: func(dir string) error {
				err := os.Mkdir(filepath.Join(dir, ".gitlab-ci-verify"), 0755)
				if err != nil {
					return err
				}
				return os.Mkdir(filepath.Join(dir, ".gitlab-ci-verify", "checks"), 0755)
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := tt.setup(dir); err != nil {
				t.Fatalf("setup failed: %s", err)
			}

			result := HasProjectPoliciesOnDisk(dir)
			if result != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
