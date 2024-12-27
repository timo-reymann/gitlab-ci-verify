package checks

import (
	_ "embed"
)

type GitlabPagesJobCheck struct {
	InMemoryCheck
}

//go:embed gitlab_pages.rego
var gitlabPagesRego string

func NewGitlabPagesJobCheck() GitlabPagesJobCheck {
	return GitlabPagesJobCheck{
		InMemoryCheck{
			RegoContent: gitlabPagesRego,
		},
	}
}
