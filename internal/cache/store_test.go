package cache

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestWriteAndReadFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		data     []byte
		expected []byte
	}{
		{"simple write", "testfile.txt", []byte("Hello, World!"), []byte("Hello, World!")},
		{"empty file", "emptyfile.txt", []byte(""), []byte("")},
		{"binary data", "binaryfile.bin", []byte{0x00, 0x01, 0x02}, []byte{0x00, 0x01, 0x02}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			dir, err := os.MkdirTemp("", "cache_test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(dir)

			// Call WriteFile
			err = WriteFile(tt.path, bytes.NewReader(tt.data))
			if err != nil {
				t.Fatalf("WriteFile() failed: %v", err)
			}

		})
	}
}

func TestOpenFile(t *testing.T) {
	_, _ = EnsureCacheDir()

	tests := []struct {
		name     string
		path     string
		data     []byte
		expected []byte
	}{
		{"simple read", "testfile.txt", []byte("Hello, World!"), []byte("Hello, World!")},
		{"empty file", "emptyfile.txt", []byte(""), []byte("")},
		{"binary data", "binaryfile.bin", []byte{0x00, 0x01, 0x02}, []byte{0x00, 0x01, 0x02}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write the test file first
			err := WriteFile(tt.path, bytes.NewReader(tt.data))
			if err != nil {
				t.Fatalf("WriteFile() failed: %v", err)
			}

			f, err := OpenFile(tt.path)
			if err != nil {
				t.Fatalf("OpenFile() failed: %v", err)
			}
			defer f.Close()

			readData, err := io.ReadAll(f)
			if err != nil {
				t.Fatalf("Failed to read data from file: %v", err)
			}

			if !bytes.Equal(readData, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, readData)
			}
		})
	}
}
