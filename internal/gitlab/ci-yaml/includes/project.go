package includes

import (
	"gopkg.in/yaml.v3"
	"slices"
)

// ProjectInclude represents a project include
// See more at https://docs.gitlab.com/ci/yaml/#includeproject
type ProjectInclude struct {
	Project string
	Files   []string
	*BaseInclude
}

func (p *ProjectInclude) Type() string {
	return "project"
}

func (p *ProjectInclude) Equals(i Include) bool {
	projectInclude, ok := i.(*ProjectInclude)
	if !ok {
		return false
	}

	if !slices.Equal(projectInclude.Files, p.Files) {
		return false
	}

	return projectInclude.Project == p.Project
}

func NewProjectInclude(node *yaml.Node, project string, files []string) *ProjectInclude {
	return &ProjectInclude{
		Project:     project,
		Files:       files,
		BaseInclude: NewBaseInclude(node),
	}
}
