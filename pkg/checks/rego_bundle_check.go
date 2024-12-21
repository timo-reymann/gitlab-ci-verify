package checks

import (
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/internal/rego_policies"
	"log"
)

type RepoBundleCheck struct {
	BundlePath string
}

func (r RepoBundleCheck) convertToCheckFinding(raw map[string]any) CheckFinding {
	cf := CheckFinding{
		Line: -1,
	}

	for k, v := range raw {
		switch k {
		case "code":
			cf.Code, _ = v.(string)
			break
		case "severity":
			severityName, _ := v.(string)
			cf.Severity = SeverityNameToLevel(severityName)
			break
		case "message":
			cf.Message, _ = v.(string)
			break
		case "link":
			cf.Link, _ = v.(string)
			break
		case "line":
			cf.Line, _ = v.(int)
			break
		default:
			logging.Debug("Ignoring field", k, "for result from bundle", r.BundlePath)
			break
		}
	}

	return cf
}

func (r RepoBundleCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	ctx := context.Background()
	rpm := rego_policies.NewRegoPolicyManager()

	if err := rpm.LoadBundle(r.BundlePath); err != nil {
		return nil, fmt.Errorf("failed to load rego bundle %s: %s", r.BundlePath, err)
	}

	regoCtx := rpm.NewRegoCtx()
	query, err := regoCtx.PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %s", err)
	}

	// TODO Add proper input
	results, err := query.Eval(ctx, rego.EvalInput(
		map[string]any{
			"yaml":       i.CiYaml.ParsedYamlMap,
			"mergedYaml": i.MergedCiYaml.ParsedYamlMap,
		},
	))
	if err != nil {
		log.Fatalf("Failed to evaluate policy: %v", err)
	}

	checkFindings := make([]CheckFinding, 0)
	for _, result := range results {
		for _, expression := range result.Expressions {
			findings := expression.Value.([]any)
			for _, finding := range findings {
				fields, ok := finding.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("finding is not of map type: %v", finding)
				}

				checkFindings = append(checkFindings, r.convertToCheckFinding(fields))
			}
		}
	}

	return checkFindings, nil
}
