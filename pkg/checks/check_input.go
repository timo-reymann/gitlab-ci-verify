package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/internal/cli"
	"github.com/timo-reymann/gitlab-ci-verify/internal/gitlab/ci-yaml"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/location"
)

// CheckInput represents the input for a check
type CheckInput struct {
	// Configuration used to execute the verifier
	Configuration *cli.Configuration
	// LintAPIResult contains the API response from the GitLab CI Lint API. It might be not set when configuration
	// disabled ci validation inside CI
	LintAPIResult *ci_yaml.VerificationResultWithRemoteInfo
	// MergedCiYaml contains the merged YAML when the lint api result is available
	MergedCiYaml *ci_yaml.CiYamlFile
	// VirtualCiYaml contains the virtual CI YAML file
	VirtualCiYaml *ci_yaml.VirtualCiYamlFile
}

// HasLintAPIResult checks if the input has a lint API result
func (c *CheckInput) HasLintAPIResult() bool {
	return c.LintAPIResult != nil
}

// CanProvideMergedYaml checks if the input can provide a merged YAML
// This is only the case when the lint API result is available
func (c *CheckInput) CanProvideMergedYaml() bool {
	return c.LintAPIResult != nil
}

// ResolveLocation resolves the location of a line in the virtual CI YAML
// Returns nil if the line is not in the virtual CI YAML
func (c *CheckInput) ResolveLocation(line int) *location.Location {
	_, loc := c.VirtualCiYaml.Resolve(line)
	return loc
}
