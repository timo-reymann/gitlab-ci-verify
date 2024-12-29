package cache

import (
	"os"
	"path/filepath"
)

// CacheDir to use for gitlab-ci-verify
func CacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "gitlab-ci-verify"), nil
}

// EnsureCacheDir to use for gitlab-ci-verify exists
func EnsureCacheDir() (string, error) {
	dir, err := CacheDir()
	if err != nil {
		return "", err
	}

	return dir, os.MkdirAll(dir, 0755)
}

func Exists(path string) bool {
	dir, err := CacheDir()
	if err != nil {
		return false
	}
	pathToCheck := filepath.Join(dir, path)
	_, err = os.Stat(pathToCheck)
	return err == nil
}
