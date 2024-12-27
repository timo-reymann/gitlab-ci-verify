package checks

import (
	"os"
	"path/filepath"
)

func projectPoliciesPath(path string) string {
	return filepath.Join(path, ".gitlab-ci-verify", "checks")
}

func HasProjectPoliciesOnDisk(path string) bool {
	stat, err := os.Stat(projectPoliciesPath(path))
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func RegisterProjectPolicies(path string) {
	bundle := projectPoliciesPath(path)
	register(BundleCheck{
		BundlePath: bundle,
	})
}
