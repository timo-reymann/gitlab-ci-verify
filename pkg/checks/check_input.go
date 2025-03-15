package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/internal/cli"
	ciyaml "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/location"
)

type CheckInput struct {
	// Configuration used to execute the verifier
	Configuration *cli.Configuration
	// LintAPIResult contains the API response from the GitLab CI Lint API. It might be not set when configuration
	// disabled ci validation inside CI
	LintAPIResult *ciyaml.VerificationResultWithRemoteInfo
	// MergedCiYaml contains the merged YAML when the lint api result is available
	MergedCiYaml *ciyaml.CiYamlFile
	// VirtualCiYaml contains the virtual CI YAML file
	VirtualCiYaml *ciyaml.VirtualCiYamlFile
}

func (c *CheckInput) HasLintAPIResult() bool {
	return c.LintAPIResult != nil
}

func (c *CheckInput) CanProvideMergedYaml() bool {
	return c.LintAPIResult != nil
}

func (c *CheckInput) ResolveLocation(line int) *location.Location {
	_, loc := c.VirtualCiYaml.Resolve(line)
	return loc
}
