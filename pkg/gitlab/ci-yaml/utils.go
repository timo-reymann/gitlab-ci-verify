package ci_yaml

import "github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"

func mustPath(path *yamlpath.Path, err error) *yamlpath.Path {
	if err != nil {
		panic(err)
	}
	return path
}
