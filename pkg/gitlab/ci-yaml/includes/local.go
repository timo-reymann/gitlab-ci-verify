package includes

import "gopkg.in/yaml.v3"

// LocalInclude represents a local include
// See more at https://docs.gitlab.com/ci/yaml/#includelocal
type LocalInclude struct {
	Path string
	*BaseInclude
}

func (l *LocalInclude) Type() string {
	return "local"
}

func (l *LocalInclude) Equals(i Include) bool {
	localInclude, ok := i.(*LocalInclude)
	if !ok {
		return false
	}

	return localInclude.Path == l.Path
}

func NewLocalInclude(node *yaml.Node, path string) *LocalInclude {
	return &LocalInclude{
		Path:        path,
		BaseInclude: NewBaseInclude(node),
	}
}
