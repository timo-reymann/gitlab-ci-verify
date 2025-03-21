package includes

import "gopkg.in/yaml.v3"

// BaseInclude is the base struct for all include types
// It contains the node of the include and is used to embed in other include types
type BaseInclude struct {
	node *yaml.Node
}

// Node returns the node of the include
func (b *BaseInclude) Node() *yaml.Node {
	return b.node
}

// NewBaseInclude creates a new base include from a yaml node
func NewBaseInclude(node *yaml.Node) *BaseInclude {
	return &BaseInclude{
		node: node,
	}
}
