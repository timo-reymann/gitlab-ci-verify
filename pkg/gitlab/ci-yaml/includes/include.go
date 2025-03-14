package includes

import (
	"github.com/timo-reymann/gitlab-ci-verify/internal/yamlpathutils"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
	"slices"
)

// Include interface for all include types
type Include interface {
	// Type of include
	Type() string
	// Node returns the YAML node of to include
	Node() *yaml.Node
	// Equals checks if two includes are equal in terms of content
	Equals(i Include) bool
}

// ParseIncludes from a yaml documents parent node
func ParseIncludes(doc *yaml.Node) ([]Include, error) {
	includePath := yamlpathutils.MustPath(yamlpath.NewPath(".include"))
	includes, err := includePath.Find(doc)
	if err != nil {
		return nil, err
	}

	parsedIncludes := make([]Include, 0)
	for _, includeNode := range includes {
		parsedIncludes = append(parsedIncludes, parseIncludeNode(includeNode)...)
	}

	return parsedIncludes, nil
}

func parseIncludeNode(includeNode *yaml.Node) []Include {
	includes := make([]Include, 0)
	switch includeNode.Tag {
	case "!!str": // simple local include
		includes = append(includes, NewLocalInclude(includeNode, includeNode.Value))

	case "!!seq": // list of includes
		items := includeNode.Content
		for _, item := range items {
			includes = append(includes, parseIncludeNode(item)...)
		}
		break

	case "!!map": // complex include
		var definition map[string]any
		err := includeNode.Decode(&definition)
		if err != nil {
			return includes
		}

		parsedComplexInclude := parseComplexInclude(includeNode, definition)
		if parsedComplexInclude != nil {
			includes = append(includes, parsedComplexInclude)
		}
		break
	}

	return includes
}

func parseComplexInclude(node *yaml.Node, definition map[string]any) Include {
	if local, isLocal := definition["local"]; isLocal {
		return NewLocalInclude(node, local.(string))
	}

	if project, isProject := definition["project"]; isProject {
		return parseProjectInclude(node, definition, project)
	}

	if component, isComponent := definition["component"]; isComponent {
		return NewComponentInclude(node, component.(string))
	}

	if template, isTemplate := definition["template"]; isTemplate {
		return NewTemplateInclude(node, template.(string))
	}

	if remote, isRemote := definition["remote"]; isRemote {
		integrity, hasIntegry := definition["integrity"]
		if hasIntegry {
			return NewRemoteInclude(node, remote.(string), integrity.(string))
		}

		return NewRemoteInclude(node, remote.(string), "")
	}

	return nil
}

func parseProjectInclude(node *yaml.Node, definition map[string]any, project any) Include {
	projectName := project.(string)

	filesNode, hasFiles := definition["file"]
	if !hasFiles {
		return NewProjectInclude(node, projectName, []string{})
	}

	files := make([]string, 0)
	switch filesNode.(type) {
	case string: // single file
		files = []string{filesNode.(string)}
		break

	case []string: // list of files
		files = filesNode.([]string)
		break

	case []any: // list of files, not correctly typed by yaml.v3
		for _, file := range filesNode.([]any) {
			if fileStr := file.(string); hasFiles {
				files = append(files, fileStr)
			}
		}
	}

	slices.Sort(files)
	return NewProjectInclude(node, projectName, files)
}
