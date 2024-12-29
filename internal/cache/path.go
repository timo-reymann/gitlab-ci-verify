package cache

import (
	"os"
	"path/filepath"
)

// CacheFolder to use for gitlab-ci-verify
func CacheFolder() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "gitlab-ci-verify"), nil
}
