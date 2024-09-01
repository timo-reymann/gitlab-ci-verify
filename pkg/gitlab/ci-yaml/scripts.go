package ci_yaml

import (
	"fmt"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

// ScriptPart represents a part of a script
type ScriptPart struct {
	// Content is the script content lines
	Content string
	// Node the script was extracted from
	Node *yaml.Node
}

func newScriptPart(n *yaml.Node) ScriptPart {
	return ScriptPart{
		Content: n.Value,
		Node:    n,
	}
}

// JobWithScripts represents the parsed scripts from a job and the name of the job
type JobWithScripts struct {
	// JobName contains the name of the job the scripts are part of
	JobName string
	// ScriptParts contains the script parts indexed by the key of the job
	ScriptParts map[string][]ScriptPart
}

func (jws *JobWithScripts) setPart(key string, parts []ScriptPart) {
	if len(parts) == 0 {
		return
	}
	jws.ScriptParts[key] = parts
}

func mustPath(path *yamlpath.Path, err error) *yamlpath.Path {
	if err != nil {
		panic(err)
	}
	return path
}

func getScriptFromKey(node *yaml.Node, key string) []ScriptPart {
	parts := make([]ScriptPart, 0)
	scriptPath := mustPath(yamlpath.NewPath(fmt.Sprintf(".%s", key)))
	scriptNodes, _ := scriptPath.Find(node)
	if scriptNodes == nil {
		return parts
	}

	for _, scriptNode := range scriptNodes {
		if scriptNode.Kind == yaml.SequenceNode {
			for _, nestedScriptNode := range scriptNode.Content {
				parts = append(parts, newScriptPart(nestedScriptNode))
			}
		} else {
			parts = append(parts, newScriptPart(scriptNode))
		}
	}

	return parts
}

var keysToExtractScriptsFrom = []string{
	"script",
	"before_script",
	"after_script",
}

// ExtractScripts from a YAML document
func ExtractScripts(doc *yaml.Node) chan JobWithScripts {
	ch := make(chan JobWithScripts)

	go func() {
		jobName := ""

		for _, node := range doc.Content[0].Content {
			// hierarchy is always scalar node (job key), mapping node (job definition)
			if node.Kind == yaml.ScalarNode {
				jobName = node.Value
				continue
			} else if node.Kind != yaml.MappingNode {
				// ignore invalid jobs
				continue
			}

			jws := JobWithScripts{
				JobName:     jobName,
				ScriptParts: map[string][]ScriptPart{},
			}

			for _, key := range keysToExtractScriptsFrom {
				jws.setPart(key, getScriptFromKey(node, key))
			}

			ch <- jws
		}
		close(ch)
	}()

	return ch
}