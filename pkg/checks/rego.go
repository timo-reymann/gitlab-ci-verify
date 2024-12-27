package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/types"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/internal/rego_policies"
	"github.com/timo-reymann/gitlab-ci-verify/internal/yamlpathutils"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
)

func queryManagerForFindings(rpm *rego_policies.RegoPolicyManager, i *CheckInput) (rego.ResultSet, error) {
	rpm.AddBuiltinFunc1("yamlPathToLineNumber", types.NewFunction(types.Args(types.S), types.N), func(context rego.BuiltinContext, param1 *ast.Term) (*ast.Term, error) {
		var yamlPathVal string
		if err := ast.As(param1.Value, &yamlPathVal); err != nil {
			return nil, err
		}

		line := yamlpathutils.PathToFirstLineNumber(i.CiYaml.ParsedYamlDoc, yamlpathutils.MustPath(yamlpath.NewPath(yamlPathVal)))
		val, err := ast.InterfaceToValue(line)
		if err != nil {
			return nil, err
		}

		return ast.NewTerm(val), nil
	})

	regoCtx := rpm.NewRegoCtx()
	ctx := context.Background()
	query, err := regoCtx.PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %s", err)
	}

	var mergedCiYaml *map[string]any

	if i.MergedCiYaml != nil {
		mergedCiYaml = &i.MergedCiYaml.ParsedYamlMap
	}

	return query.Eval(ctx, rego.EvalInput(
		map[string]any{
			"yaml":       &i.CiYaml.ParsedYamlMap,
			"mergedYaml": mergedCiYaml,
		},
	))
}

func convertToCheckFinding(sourcePath string, raw map[string]any) (*CheckFinding, error) {
	convertErrs := make([]error, 0)

	addErrIfNotOk := func(ok bool, field string, expectedType string) {
		if !ok {
			convertErrs = append(convertErrs, fmt.Errorf("failed to convert field %s to %s", field, expectedType))
		}
	}

	cf := CheckFinding{
		Line:     -1,
		Severity: -1,
	}
	ok := false

	for k, v := range raw {
		switch k {
		case "code":
			cf.Code, ok = v.(string)
			addErrIfNotOk(ok, "code", "string")
			break

		case "severity":
			severityName, ok := v.(string)
			addErrIfNotOk(ok, "severity", "string")
			cf.Severity = SeverityNameToLevel(severityName)
			break

		case "message":
			cf.Message, ok = v.(string)
			addErrIfNotOk(ok, "message", "string")
			break

		case "link":
			cf.Link, ok = v.(string)
			addErrIfNotOk(ok, "link", "string")
			break

		case "line":
			num, ok := v.(json.Number)
			addErrIfNotOk(ok, "line", "number")
			numInt, err := num.Int64()
			addErrIfNotOk(err == nil, "line", "int64")
			cf.Line = int(numInt)
			break

		default:
			logging.Debug("Ignoring field", k, "for result from", sourcePath)
			break
		}
	}

	if cf.Code == "" {
		convertErrs = append(convertErrs, errors.New("code is required"))
	}

	if cf.Message == "" {
		convertErrs = append(convertErrs, errors.New("message is required"))
	}

	if cf.Severity == -1 {
		convertErrs = append(convertErrs, errors.New("severity is required"))
	}

	return &cf, errors.Join(convertErrs...)
}

func parseResults(regoPath string, results rego.ResultSet) ([]CheckFinding, error) {
	parseErrs := make([]error, 0)
	checkFindings := make([]CheckFinding, 0)

	logging.Verbose("Got rego results", results, "for", regoPath)
	for _, result := range results {
		for _, expression := range result.Expressions {
			findings := expression.Value.([]any)
			for _, finding := range findings {
				fields, ok := finding.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("finding is not of map type: %v", finding)
				}

				checkFinding, err := convertToCheckFinding(regoPath, fields)
				if err != nil {
					parseErrs = append(parseErrs, err)
				}
				if checkFinding != nil {
					checkFindings = append(checkFindings, *checkFinding)
				}
			}
		}
	}

	var err error
	if len(parseErrs) > 0 {
		err = errors.Join(parseErrs...)
	} else {
		err = nil
	}

	return checkFindings, err
}
