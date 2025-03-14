package yamlpathutils

import (
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

// PathToLineNumbers runs the path on the given node and returns all node line numbers that resolved from the query
func PathToLineNumbers(node *yaml.Node, path *yamlpath.Path) []int {
	nodes, err := path.Find(node)
	if err != nil {
		return []int{}
	}

	lines := make([]int, len(nodes))
	for idx, n := range nodes {
		lines[idx] = n.Line
	}

	return lines
}

// PathToFirstLineNumber runs the path on the given node and returns the first line number that resolved from the query
func PathToFirstLineNumber(node *yaml.Node, path *yamlpath.Path) int {
	lines := PathToLineNumbers(node, path)
	if len(lines) > 0 {
		return lines[0]
	}
	return -1
}
