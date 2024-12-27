package checks

import (
	"encoding/json"
	"github.com/timo-reymann/gitlab-ci-verify/internal/rego_policies"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	"testing"

	"github.com/open-policy-agent/opa/v1/rego"
)

func TestParseResults(t *testing.T) {
	testCases := []struct {
		name             string
		results          rego.ResultSet
		expectedFindings []CheckFinding
		expectError      bool
	}{
		{
			name: "Single valid finding",
			results: rego.ResultSet{
				{
					Expressions: []*rego.ExpressionValue{
						{
							Value: []any{
								map[string]any{
									"code":     "TEST_CODE",
									"severity": "error",
									"message":  "Test message",
									"link":     "http://example.com",
									"line":     json.Number("42"),
								},
							},
						},
					},
				},
			},
			expectedFindings: []CheckFinding{
				{
					Code:     "TEST_CODE",
					Severity: SeverityError,
					Message:  "Test message",
					Link:     "http://example.com",
					Line:     42,
				},
			},
			expectError: false,
		},
		{
			name: "Extra fields",
			results: rego.ResultSet{
				{
					Expressions: []*rego.ExpressionValue{
						{
							Value: []any{
								map[string]any{
									"code":      "TEST_CODE",
									"severity":  "error",
									"message":   "Test message",
									"link":      "http://example.com",
									"line":      json.Number("42"),
									"ignore_me": true,
								},
							},
						},
					},
				},
			},
			expectedFindings: []CheckFinding{
				{
					Code:     "TEST_CODE",
					Severity: SeverityError,
					Message:  "Test message",
					Link:     "http://example.com",
					Line:     42,
				},
			},
			expectError: false,
		},
		{
			name: "Missing code",
			results: rego.ResultSet{
				{
					Expressions: []*rego.ExpressionValue{
						{
							Value: []any{
								map[string]any{
									"severity": "error",
									"message":  "Test message",
									"link":     "http://example.com",
									"line":     json.Number("42"),
								},
							},
						},
					},
				},
			},
			expectedFindings: nil,
			expectError:      true,
		},
		{
			name: "Missing severity",
			results: rego.ResultSet{
				{
					Expressions: []*rego.ExpressionValue{
						{
							Value: []any{
								map[string]any{
									"code":    "TEST_CODE",
									"message": "Test message",
									"link":    "http://example.com",
									"line":    json.Number("42"),
								},
							},
						},
					},
				},
			},
			expectedFindings: nil,
			expectError:      true,
		},
		{
			name: "Missing message",
			results: rego.ResultSet{
				{
					Expressions: []*rego.ExpressionValue{
						{
							Value: []any{
								map[string]any{
									"code":     "TEST_CODE",
									"severity": "error",
									"link":     "http://example.com",
									"line":     json.Number("42"),
								},
							},
						},
					},
				},
			},
			expectedFindings: nil,
			expectError:      true,
		},
		{
			name: "Invalid line type",
			results: rego.ResultSet{
				{
					Expressions: []*rego.ExpressionValue{
						{
							Value: []any{
								map[string]any{
									"code":     "TEST_CODE",
									"severity": "error",
									"message":  "Test message",
									"link":     "http://example.com",
									"line":     "invalid_line",
								},
							},
						},
					},
				},
			},
			expectedFindings: nil,
			expectError:      true,
		},
		{
			name: "Invalid finding type",
			results: rego.ResultSet{
				{
					Expressions: []*rego.ExpressionValue{
						{
							Value: []any{
								"invalid_finding",
							},
						},
					},
				},
			},
			expectedFindings: nil,
			expectError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkFindings, err := parseResults("test.rego", tc.results)
			if (err != nil) != tc.expectError {
				t.Fatalf("parseResults() error = %v, expectError %v", err, tc.expectError)
			}
			if err == nil && len(checkFindings) != len(tc.expectedFindings) {
				t.Fatalf("expected %d findings, got %d", len(tc.expectedFindings), len(checkFindings))
			}
			for i, expectedFinding := range tc.expectedFindings {
				if checkFindings[i] != expectedFinding {
					t.Fatalf("expected finding %+v, got %+v", expectedFinding, checkFindings[i])
				}
			}
		})
	}
}

func TestQueryRegoForFindings(t *testing.T) {
	rpm := rego_policies.NewRegoPolicyManager()
	err := rpm.LoadBundle("test_data/rego-bundle")
	if err != nil {
		t.Fatalf("failed to load rego bundle: %s", err)
	}

	input := &CheckInput{
		CiYaml:        NewCiYamlFromFile(t, "test_data/ci_yamls/singleBuild.yml"),
		MergedCiYaml:  NewCiYamlFromFile(t, "test_data/ci_yamls/singleBuild.yml"),
		Configuration: &cli.Configuration{},
	}

	results, err := queryManagerForFindings(rpm, input)
	if err != nil {
		t.Fatalf("queryManagerForFindings failed: %s", err)
	}

	if len(results) == 0 {
		t.Fatalf("expected findings, got none")
	}
}
