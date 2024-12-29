package cache

import (
	"io"
	"os"
	"path/filepath"
)

func WriteFile(path string, r io.Reader) error {
	dir, err := EnsureCacheDir()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(dir, path), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	return err
}

func OpenFile(path string) (*os.File, error) {
	dir, err := EnsureCacheDir()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath.Join(dir, path))
	return f, err
}
