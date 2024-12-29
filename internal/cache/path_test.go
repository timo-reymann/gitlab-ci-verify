package cache

import "testing"

func TestCacheFolder(t *testing.T) {
	_, err := CacheFolder()
	if err != nil {
		t.Errorf("CacheFolder() failed: %v", err)
	}
}
