package checks

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/rego_policies"
)

type RemoteBundleCheck struct {
	BundleURL string
}

func (r RemoteBundleCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	rpm := rego_policies.NewRegoPolicyManager()

	logging.Verbose("Load remote rego bundle", r.BundleURL)
	if err := rpm.LoadBundleFromRemote(r.BundleURL); err != nil {
		return nil, fmt.Errorf("failed to load rego bundle %s: %s", r.BundleURL, err)
	}

	results, err := queryManagerForFindings(rpm, i)
	if err != nil {
		logging.Verbose("Failed to query rego remote bundle", r.BundleURL, err)
		return nil, err
	}

	return parseResults(i, r.BundleURL, results)
}
