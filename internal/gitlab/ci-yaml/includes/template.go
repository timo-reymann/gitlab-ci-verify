package includes

import "gopkg.in/yaml.v3"

// TemplateInclude represents a template include
// See more at https://docs.gitlab.com/ci/yaml/#includetemplate
type TemplateInclude struct {
	*BaseInclude
	Template string
}

func (t *TemplateInclude) Type() string {
	return "template"
}

func (t *TemplateInclude) Equals(i Include) bool {
	templateInclude, ok := i.(*TemplateInclude)
	if !ok {
		return false
	}

	return templateInclude.Template == t.Template
}

func NewTemplateInclude(node *yaml.Node, template string) *TemplateInclude {
	return &TemplateInclude{
		Template:    template,
		BaseInclude: NewBaseInclude(node),
	}
}
