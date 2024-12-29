package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCacheFolder(t *testing.T) {
	_, err := CacheDir()
	if err != nil {
		t.Errorf("CacheDir() failed: %v", err)
	}
}

func TestExists(t *testing.T) {
	dir, err := CacheDir()
	if err != nil {
		t.Fatalf("Failed to get cache dir: %v", err)
	}

	// Test matrix
	tests := []struct {
		name     string
		filename string
		exists   bool
	}{
		{"File exists", "testfile", true},
		{"File does not exist", "nonexistent", false},
	}

	cacheDir, err := CacheDir()
	if err != nil {
		t.Errorf("CacheDir() failed: %v", err)
	}
	os.RemoveAll(cacheDir)
	_, _ = EnsureCacheDir()

	// Create a file in the temporary directory
	testFile := filepath.Join(dir, "testfile")
	_ = os.Remove(testFile)
	if _, err := os.Create(testFile); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Exists(tt.filename) != tt.exists {
				t.Errorf("Exists(%q) = %v, want %v", tt.filename, !tt.exists, tt.exists)
			}
		})
	}
}
