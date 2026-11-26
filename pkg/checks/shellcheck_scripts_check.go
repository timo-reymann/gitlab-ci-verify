package checks

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"sync"

	ci_yaml2 "github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/ci-yaml"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/shellcheck"
	"gopkg.in/yaml.v3"
)

type ShellScriptCheck struct {
}

// analyzedScript is a script of a job that was analyzed by shellcheck
type analyzedScript struct {
	// JobName contains the name of the job the script is part of
	JobName string
	// Key is the key of the script inside the job, e.g., before_script
	Key string
	// Lines contains the lines of the script
	Lines []ci_yaml2.ScriptPartLine
	// Script contains the joined script that was analyzed
	Script []byte
	// Findings from shellcheck
	Findings []shellcheck.Finding
}

// line returns the script line for the 1-based shellcheck line
func (a *analyzedScript) line(shellcheckLine int) ci_yaml2.ScriptPartLine {
	// handle cases where the YAML parsing does not match the actual
	// line, e.g., with multiline list continuation
	lineNo := shellcheckLine - 1
	if lineNo < 0 || lineNo >= len(a.Lines) {
		lineNo = 0
	}
	return a.Lines[lineNo]
}

func scriptKey(jobName, key string) string {
	return jobName + "\x00" + key
}

func (s ShellScriptCheck) shellcheckLevelToSeverity(level string) int {
	switch level {
	case "error":
		return SeverityError
	case "warning":
		return SeverityWarning
	case "info":
		return SeverityInfo
	case "style":
		return SeverityStyle
	}
	return -1
}

func (s ShellScriptCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	shellChecker, err := shellcheck.NewShellChecker()
	if err != nil {
		return nil, err
	}
	defer shellChecker.Close()

	localScripts := s.analyze(shellChecker, i.Configuration.ShellcheckFlags, i.VirtualCiYaml.Combined.ParsedYamlDoc, nil)
	findings := s.virtualFileFindings(i, localScripts)

	if i.Configuration.ShellcheckMerged {
		if i.CanProvideMergedYaml() {
			findings = append(findings, s.mergedYamlFindings(shellChecker, i, localScripts)...)
		} else {
			logging.Warn("Shellcheck for the merged YAML is enabled, but the merged YAML is not available as the CI Lint API is not used or the pipeline is invalid")
		}
	}

	return findings, nil
}

// analyze all scripts in the given document, scripts for which skip returns true are not analyzed
func (s ShellScriptCheck) analyze(shellChecker *shellcheck.ShellChecker, flags string, doc *yaml.Node, skip func(jobName, key string, script []byte) bool) []analyzedScript {
	var mu sync.Mutex
	var wg sync.WaitGroup
	// results are indexed by the position of the job, to keep the order stable
	resultsByJob := map[int][]analyzedScript{}
	jobCount := 0

	for jobWithScripts := range ci_yaml2.ExtractScripts(doc) {
		jobIdx := jobCount
		jobCount++
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, key := range slices.Sorted(maps.Keys(jobWithScripts.ScriptParts)) {
				lines, joinedScript := ci_yaml2.Concat(jobWithScripts.ScriptParts[key])
				if skip != nil && skip(jobWithScripts.JobName, key, joinedScript) {
					continue
				}

				result, err := shellChecker.AnalyzeSnippet(joinedScript, flags)
				if err != nil {
					logging.Warn("Failed to analyze snippet in job", jobWithScripts.JobName)
					continue
				}

				mu.Lock()
				resultsByJob[jobIdx] = append(resultsByJob[jobIdx], analyzedScript{
					JobName:  jobWithScripts.JobName,
					Key:      key,
					Lines:    lines,
					Script:   joinedScript,
					Findings: result.Findings,
				})
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	results := make([]analyzedScript, 0, jobCount)
	for idx := 0; idx < jobCount; idx++ {
		results = append(results, resultsByJob[idx]...)
	}
	return results
}

func (s ShellScriptCheck) virtualFileFindings(i *CheckInput, scripts []analyzedScript) []CheckFinding {
	findings := make([]CheckFinding, 0)
	for _, script := range scripts {
		for _, f := range script.Findings {
			line := script.line(f.Line)
			loc := i.ResolveLocation(line.Node.Line)
			if loc == nil {
				panic("Failed to map location for shellcheck finding in line " + strconv.Itoa(line.Node.Line))
			}

			findings = append(findings, CheckFinding{
				Severity: s.shellcheckLevelToSeverity(f.Level),
				Code:     fmt.Sprintf("SC-%s", strconv.Itoa(f.Code)),
				Line:     loc.Line + 1,
				Message:  fmt.Sprintf("[%s:%s:%d] %s", script.JobName, script.Key, f.Line, f.Message),
				Link:     fmt.Sprintf("https://www.shellcheck.net/wiki/SC%d", f.Code),
				File:     loc.File,
			})
		}
	}
	return findings
}

// mergedYamlFindings analyzes the scripts of the merged YAML that differ from the scripts in the virtual CI YAML,
// e.g., scripts from CI/CD components or scripts with interpolated inputs.
//
// As the merged YAML contains no information about the origin of a script, findings are reported for the job in
// the virtual CI YAML if it exists, and the entry file otherwise.
func (s ShellScriptCheck) mergedYamlFindings(shellChecker *shellcheck.ShellChecker, i *CheckInput, localScripts []analyzedScript) []CheckFinding {
	knownScripts := map[string]struct{}{}
	knownLines := map[string]map[string]struct{}{}
	for _, script := range localScripts {
		knownScripts[string(script.Script)] = struct{}{}

		lineSet := map[string]struct{}{}
		for _, l := range script.Lines {
			lineSet[strings.TrimSpace(l.LineContent)] = struct{}{}
		}
		knownLines[scriptKey(script.JobName, script.Key)] = lineSet
	}

	mergedScripts := s.analyze(shellChecker, i.Configuration.ShellcheckFlags, i.MergedCiYaml.ParsedYamlDoc, func(_, _ string, script []byte) bool {
		// identical scripts are already analyzed as part of the virtual CI YAML
		_, known := knownScripts[string(script)]
		return known
	})

	jobLines := topLevelKeyLines(i.VirtualCiYaml.Combined.ParsedYamlDoc)
	findings := make([]CheckFinding, 0)
	for _, script := range mergedScripts {
		lineSet := knownLines[scriptKey(script.JobName, script.Key)]
		file, line := i.VirtualCiYaml.EntryFilePath, -1
		if jobLine, ok := jobLines[script.JobName]; ok {
			if loc := i.ResolveLocation(jobLine); loc != nil {
				file, line = loc.File, loc.Line+1
			}
		}

		for _, f := range script.Findings {
			scriptLine := strings.TrimSpace(script.line(f.Line).LineContent)
			// unchanged lines of a job are already analyzed as part of the virtual CI YAML
			if _, known := lineSet[scriptLine]; known {
				continue
			}

			findings = append(findings, CheckFinding{
				Severity: s.shellcheckLevelToSeverity(f.Level),
				Code:     fmt.Sprintf("SC-%s", strconv.Itoa(f.Code)),
				Line:     line,
				Message:  fmt.Sprintf("[merged:%s:%s:%d] %s (in merged YAML: %s)", script.JobName, script.Key, f.Line, f.Message, truncate(scriptLine, 80)),
				Link:     fmt.Sprintf("https://www.shellcheck.net/wiki/SC%d", f.Code),
				File:     file,
			})
		}
	}

	return findings
}

// topLevelKeyLines returns the line for all top level keys of the document
func topLevelKeyLines(doc *yaml.Node) map[string]int {
	lines := map[string]int{}
	if doc == nil || len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return lines
	}

	root := doc.Content[0]
	for idx := 0; idx < len(root.Content); idx += 2 {
		if _, exists := lines[root.Content[idx].Value]; !exists {
			lines[root.Content[idx].Value] = root.Content[idx].Line
		}
	}
	return lines
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
