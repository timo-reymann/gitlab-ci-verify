package checks

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/internal/rego_policies"
)

type BundleCheck struct {
	BundlePath string
}

func (r BundleCheck) convertToCheckFinding(raw map[string]any) (*CheckFinding, error) {
	return convertToCheckFinding(r.BundlePath, raw)
}

func (r BundleCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	rpm := rego_policies.NewRegoPolicyManager()

	logging.Verbose("Load rego bundle", r.BundlePath)
	if err := rpm.LoadBundle(r.BundlePath); err != nil {
		return nil, fmt.Errorf("failed to load rego bundle %s: %s", r.BundlePath, err)
	}

	results, err := queryManagerForFindings(rpm, i)
	if err != nil {
		logging.Verbose("Failed to query rego bundle", r.BundlePath, err)
		return nil, err
	}

	return parseResults(r.BundlePath, results)
}
