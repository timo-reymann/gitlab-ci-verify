package yamlpathutils

import (
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
)

// MustPath ensures that a path is valid and if it is invalid panics
func MustPath(path *yamlpath.Path, err error) *yamlpath.Path {
	if err != nil {
		panic(err)
	}
	return path
}
