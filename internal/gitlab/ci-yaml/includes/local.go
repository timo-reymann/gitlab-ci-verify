package includes

import (
	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
	"path"
	"sort"
	"strings"
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
	return path.Join(projectDir, srcDir, l.Path)
}

// IsGlobPattern returns true if the path contains glob pattern characters
func (l *LocalInclude) IsGlobPattern() bool {
	return strings.ContainsAny(l.Path, "*?[") || strings.Contains(l.Path, "**")
}

// ResolvePaths returns all matching paths, expanding globs if present
// If the path is not a glob pattern, it returns a slice containing the single resolved path
// If the path is a glob pattern, it expands the pattern and returns all matching paths
// Empty matches return an empty slice (GitLab silently skips these)
func (l *LocalInclude) ResolvePaths(projectDir, srcFile string) ([]string, error) {
	basePath := l.ResolvePath(projectDir, srcFile)

	if !l.IsGlobPattern() {
		return []string{basePath}, nil
	}

	matches, err := doublestar.FilepathGlob(basePath)
	if err != nil {
		return nil, err
	}

	// Empty match is valid (GitLab silently skips)
	if len(matches) == 0 {
		return []string{}, nil
	}

	sort.Strings(matches) // Deterministic ordering
	return matches, nil
}

func NewLocalInclude(node *yaml.Node, path string) *LocalInclude {
	return &LocalInclude{
		Path:        path,
		BaseInclude: NewBaseInclude(node),
	}
}
