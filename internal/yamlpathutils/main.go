package yamlpathutils

import "github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"

func MustPath(path *yamlpath.Path, err error) *yamlpath.Path {
	if err != nil {
		panic(err)
	}
	return path
}
