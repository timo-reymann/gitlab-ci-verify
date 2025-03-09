package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	ciyaml "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
)

type CheckInput struct {
	// CiYaml contains the YAML configuration in different representations as loaded from the file system
	CiYaml *ciyaml.CiYamlFile
	// Configuration used to execute the verifier
	Configuration *cli.Configuration
	// LintAPIResult contains the API response from the GitLab CI Lint API. It might be not set when configuration
	// disabled ci validation inside CI
	LintAPIResult *ciyaml.VerificationResultWithRemoteInfo
	// MergedCiYaml contains the merged YAML when the lint api result is available
	MergedCiYaml *ciyaml.CiYamlFile
}

func (c *CheckInput) HasLintAPIResult() bool {
	return c.LintAPIResult != nil
}

func (c *CheckInput) CanProvideMergedYaml() bool {
	return c.LintAPIResult != nil
}
