package cli

import (
	"bytes"
	"os"
	"testing"
)

func TestReadStdinPipe(t *testing.T) {
	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	testData := []byte("test input")
	_, err = w.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write to pipe: %v", err)
	}
	_ = w.Close()
	os.Stdin = r

	result, err := ReadStdinPipe()
	if err != nil {
		t.Fatalf("ReadStdinPipe() error = %v, want nil", err)
	}

	if !bytes.Equal(result, testData) {
		t.Errorf("ReadStdinPipe() = %v, want %v", result, testData)
	}
}
