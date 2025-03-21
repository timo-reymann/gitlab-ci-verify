package includes

import "gopkg.in/yaml.v3"

// ComponentInclude represents a component include
// See more at https://docs.gitlab.com/ci/yaml/#includecomponent
type ComponentInclude struct {
	*BaseInclude
	Component string
}

func (c *ComponentInclude) Type() string {
	return "component"
}

func (c *ComponentInclude) Equals(i Include) bool {
	componentInclude, ok := i.(*ComponentInclude)
	if !ok {
		return false
	}

	return componentInclude.Component == c.Component
}

func NewComponentInclude(node *yaml.Node, component string) *ComponentInclude {
	return &ComponentInclude{
		Component:   component,
		BaseInclude: NewBaseInclude(node),
	}
}
