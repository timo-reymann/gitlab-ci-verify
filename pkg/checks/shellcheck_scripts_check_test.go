package checks

import (
	"path"
	"slices"
	"testing"
)

func TestShellScriptCheck_Run(t *testing.T) {
	c := ShellScriptCheck{}
	testCases := []struct {
		name             string
		file             string
		expectedFindings []CheckFinding
	}{
		{
			name: "With info finding",
			file: "withScriptInfoFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityInfo,
					Code:     "SC-2086",
					Line:     4,
					Message:  "[build:script:1] Double quote to prevent globbing and word splitting.",
					Link:     "https://www.shellcheck.net/wiki/SC2086",
					File:     "withScriptInfoFinding.yml",
				},
			},
		},
		{
			name: "With error finding",
			file: "withScriptErrorFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityError,
					Code:     "SC-2189",
					Line:     4,
					Message:  "[build:script:1] You can't have | between this redirection and the command it should apply to.",
					Link:     "https://www.shellcheck.net/wiki/SC2189",
					File:     "withScriptErrorFinding.yml",
				},
			},
		},
		{
			name: "With warning finding",
			file: "withScriptWarningFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityWarning,
					Code:     "SC-2202",
					Line:     4,
					Message:  "[build:script:1] Globs don't work as operands in [ ]. Use a loop.",
					Link:     "https://www.shellcheck.net/wiki/SC2202",
					File:     "withScriptWarningFinding.yml",
				},
			},
		},
		{
			name: "With style finding",
			file: "withStyleWarningFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityStyle,
					Code:     "SC-2006",
					Line:     4,
					Message:  "[build:script:1] Use $(...) notation instead of legacy backticks `...`.",
					Link:     "https://www.shellcheck.net/wiki/SC2006",
					File:     "withStyleWarningFinding.yml",
				},
			},
		},
		{
			name:             "With no findings and reference",
			file:             "withReferenceNoFinding.yml",
			expectedFindings: []CheckFinding{},
		},
		{
			name: "With script list item continuation",
			file: "withScriptListItemContinuation.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityError,
					Code:     "SC-1070",
					Line:     4,
					Message:  `[build:script:1] Parsing stopped here. Mismatched keywords or invalid parentheses?`,
					Link:     "https://www.shellcheck.net/wiki/SC1070",
					File:     "withScriptListItemContinuation.yml",
				},
				{
					Severity: SeverityError,
					Code:     "SC-1141",
					Line:     4,
					Message:  `[build:script:1] Unexpected tokens after compound command. Bad redirection or missing ;/&&/||/|?`,
					Link:     "https://www.shellcheck.net/wiki/SC1141",
					File:     "withScriptListItemContinuation.yml",
				},
			},
		},
		{
			name: "Nested includes with findings",
			file: "withNestedIncludesAndFindings.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityInfo,
					Code:     "SC-2086",
					Line:     3,
					Message:  "[.mod_download:before_script:1] Double quote to prevent globbing and word splitting.",
					Link:     "https://www.shellcheck.net/wiki/SC2086",
					File:     "test_data/ci_yamls/includes/ci/templates/.mod_download.gitlab-ci.yml",
				},
			},
		},
		{
			name: "With spec header and finding",
			file: "withSpecHeaderAndFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityInfo,
					Code:     "SC-2086",
					Line:     9,
					Message:  "[build:script:1] Double quote to prevent globbing and word splitting.",
					Link:     "https://www.shellcheck.net/wiki/SC2086",
					File:     "withSpecHeaderAndFinding.yml",
				},
			},
		},
		{
			name: "Local include with spec header and finding",
			file: "withLocalComponentIncludeAndFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityInfo,
					Code:     "SC-2086",
					Line:     9,
					Message:  "[test:script:1] Double quote to prevent globbing and word splitting.",
					Link:     "https://www.shellcheck.net/wiki/SC2086",
					File:     "test_data/ci_yamls/includes/ci/local-components/test.yml",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkInput := createCheckInput(t, NewCiYamlFromFile(t, path.Join("test_data", "ci_yamls", tc.file)), path.Join("test_data", "ci_yamls"), tc.file)
			VerifyFindings(t, tc.expectedFindings, CheckMustSucceed(c.Run(checkInput)))
		})
	}
}

func TestShellScriptCheck_RunMergedYaml(t *testing.T) {
	const file = "withRemoteComponentForMergedYaml.yml"
	localFindings := []CheckFinding{
		{
			Severity: SeverityInfo,
			Code:     "SC-2086",
			Line:     7,
			Message:  "[local-job:script:2] Double quote to prevent globbing and word splitting.",
			Link:     "https://www.shellcheck.net/wiki/SC2086",
			File:     file,
		},
		{
			Severity: SeverityInfo,
			Code:     "SC-2086",
			Line:     11,
			Message:  "[partial-job:script:1] Double quote to prevent globbing and word splitting.",
			Link:     "https://www.shellcheck.net/wiki/SC2086",
			File:     file,
		},
	}
	mergedFindings := []CheckFinding{
		{
			Severity: SeverityInfo,
			Code:     "SC-2086",
			Line:     9,
			Message:  "[merged:partial-job:script:2] Double quote to prevent globbing and word splitting. (in merged YAML: echo $INTERPOLATED)",
			Link:     "https://www.shellcheck.net/wiki/SC2086",
			File:     file,
		},
		{
			Severity: SeverityInfo,
			Code:     "SC-2086",
			Line:     -1,
			Message:  "[merged:component-job:script:1] Double quote to prevent globbing and word splitting. (in merged YAML: echo $UNQUOTED_COMPONENT)",
			Link:     "https://www.shellcheck.net/wiki/SC2086",
			File:     file,
		},
	}

	testCases := []struct {
		name             string
		enabled          bool
		withMergedYaml   bool
		expectedFindings []CheckFinding
	}{
		{
			name:             "disabled",
			enabled:          false,
			withMergedYaml:   true,
			expectedFindings: localFindings,
		},
		{
			name:             "enabled without merged YAML",
			enabled:          true,
			withMergedYaml:   false,
			expectedFindings: localFindings,
		},
		{
			name:             "enabled",
			enabled:          true,
			withMergedYaml:   true,
			expectedFindings: append(slices.Clone(localFindings), mergedFindings...),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkInput := createCheckInput(t, NewCiYamlFromFile(t, path.Join("test_data", "ci_yamls", file)), path.Join("test_data", "ci_yamls"), file)
			checkInput.Configuration.ShellcheckMerged = tc.enabled
			checkInput.MergedCiYaml = nil
			if tc.withMergedYaml {
				checkInput.MergedCiYaml = NewCiYamlFromFile(t, path.Join("test_data", "ci_yamls", "withRemoteComponentForMergedYaml.merged.yml"))
			}

			VerifyFindings(t, tc.expectedFindings, CheckMustSucceed(ShellScriptCheck{}.Run(checkInput)))
		})
	}
}
