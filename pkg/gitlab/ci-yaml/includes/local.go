package includes

import (
	"gopkg.in/yaml.v3"
	"path"
)

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

// ResolvePath resolves the path of the local include
// If the path is absolute, it is resolved from the project root
// If the path is relative, it is resolved from the source file that includes it
func (l *LocalInclude) ResolvePath(projectDir, srcFile string) string {
	if len(l.Path) < 1 {
		return ""
	}

	// Absolute paths should always be resolved from the project root
	if path.IsAbs(l.Path) {
		return path.Join(projectDir, l.Path[1:])
	}

	// Relative paths should be resolved from the source file that includes them
	srcDir := path.Dir(srcFile)
	return path.Join(srcDir, l.Path)
}

func NewLocalInclude(node *yaml.Node, path string) *LocalInclude {
	return &LocalInclude{
		Path:        path,
		BaseInclude: NewBaseInclude(node),
	}
}
