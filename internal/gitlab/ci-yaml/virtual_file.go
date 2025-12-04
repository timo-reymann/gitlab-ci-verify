package ci_yaml

import (
	"bytes"
	includes2 "github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/ci-yaml/includes"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/filtering"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/location"
	"gopkg.in/yaml.v3"
	"os"
	"slices"
	"sort"
)

// VirtualCiYamlFilePart represents a part of a virtual CI YAML file
type VirtualCiYamlFilePart struct {
	// CiYamlFile contains the parsed CI YAML file
	CiYamlFile *CiYamlFile
	// Path of the part relative to the project root
	Path string
	// StartLine of the part (1-based)
	StartLine int
	// EndLine of the part (1-based)
	EndLine int
}

// VirtualCiYamlFile represents a virtual CI YAML file
type VirtualCiYamlFile struct {
	// EntryFile is the entry file gitlab-ci-verify was initially executed for
	EntryFile *CiYamlFile
	// EntryFilePath is the path of the file gitlab-ci-verify was initially executed for
	EntryFilePath string
	// Combined contains the concat CI YAML file
	Combined *CiYamlFile
	// Parts contains all parts of the virtual CI YAML file in the order they appear
	Parts []*VirtualCiYamlFilePart
	// ResolveFindings contains findings generated during virtual file resolution (e.g., empty glob patterns)
	ResolveFindings []VirtualFileResolveFinding
}

// VirtualFileResolveFinding represents a finding generated during virtual file resolution
type VirtualFileResolveFinding struct {
	// Code is the finding code (e.g., 101, 102)
	Code int
	// Severity is the severity level (0 = Error, 1 = Warning)
	Severity int
	// Message is the finding message
	Message string
	// IncludePath is the path of the include that generated the finding
	IncludePath string
}

// Resolve the part and location of a line in the virtual CI YAML
// Returns nil if the line is not in the virtual CI YAML
func (v *VirtualCiYamlFile) Resolve(line int) (*VirtualCiYamlFilePart, *location.Location) {
	idx := sort.Search(len(v.Parts), func(i int) bool {
		return v.Parts[i].StartLine >= line
	})
	var part *VirtualCiYamlFilePart
	if idx > 0 {
		part = v.Parts[idx-1]
	}

	if part != nil && part.StartLine <= line && part.EndLine >= line {
		return part, location.NewLocation(part.Path, line-part.StartLine)
	}

	return nil, nil
}

// GetIgnoredCodes for a line in the virtual CI YAML
// This takes into account the file level, and line level ignores
func (v *VirtualCiYamlFile) GetIgnoredCodes(line int) []string {
	ignored := make([]string, 0)

	part, _ := v.Resolve(line)
	if part == nil {
		return ignored
	}

	partCiYaml := part.CiYamlFile

	ignored = append(ignored, filtering.IgnoreCommentsToCodes(partCiYaml.GetFileLevelIgnores())...)
	ignored = append(ignored, filtering.IgnoreCommentsToCodes(partCiYaml.GetLineLevelIgnores(line))...)

	return ignored
}

// CreateVirtualCiYamlFile creates a virtual CI YAML file from a project root, entry file path, and parsed entry file
func CreateVirtualCiYamlFile(projectRoot string, entryFilePath string, entryFile *CiYamlFile) (*VirtualCiYamlFile, error) {
	virtualFile := &VirtualCiYamlFile{
		EntryFile:       entryFile,
		EntryFilePath:   entryFilePath,
		Parts:           []*VirtualCiYamlFilePart{},
		ResolveFindings: []VirtualFileResolveFinding{},
	}

	combined := bytes.NewBuffer(nil)
	line := addEntryFilePart(virtualFile, entryFile, entryFilePath, combined)

	uniqueLocalIncludes := getUniqueLocalIncludes(entryFile)
	line, err := processLocalIncludes(virtualFile, uniqueLocalIncludes, projectRoot, entryFilePath, combined, line)
	if err != nil {
		return nil, err
	}

	if err := finalizeCombinedYaml(virtualFile, uniqueLocalIncludes, combined); err != nil {
		return nil, err
	}

	return virtualFile, nil
}

// addEntryFilePart adds the entry file as the first part and returns the next line number
func addEntryFilePart(virtualFile *VirtualCiYamlFile, entryFile *CiYamlFile, entryFilePath string, combined *bytes.Buffer) int {
	entryLineCount := bytes.Count(entryFile.FileContent, []byte("\n"))
	virtualFile.Parts = append(virtualFile.Parts, &VirtualCiYamlFilePart{
		CiYamlFile: entryFile,
		Path:       entryFilePath,
		StartLine:  1,
		EndLine:    entryLineCount + 2,
	})
	combined.Write(entryFile.FileContent)
	combined.Write([]byte("\n"))
	return entryLineCount + 2
}

// getUniqueLocalIncludes extracts and deduplicates local includes from the entry file
func getUniqueLocalIncludes(entryFile *CiYamlFile) []includes2.Include {
	localIncludes := includes2.FilterByTypes(entryFile.Includes, "local")
	return includes2.Unique(localIncludes)
}

// processLocalIncludes processes all local includes and adds them to the virtual file
func processLocalIncludes(virtualFile *VirtualCiYamlFile, uniqueLocalIncludes []includes2.Include, projectRoot, entryFilePath string, combined *bytes.Buffer, startLine int) (int, error) {
	addedIncludePaths := make([]string, 0)
	line := startLine

	for _, uniqueLocalInclude := range uniqueLocalIncludes {
		localInclude := uniqueLocalInclude.(*includes2.LocalInclude)

		// Check for unsupported glob characters first
		if localInclude.HasUnsupportedGlobChars() {
			virtualFile.ResolveFindings = append(virtualFile.ResolveFindings, VirtualFileResolveFinding{
				Code:        103,
				Severity:    0, // Error level
				Message:     "Include pattern '" + localInclude.Path + "' uses unsupported glob syntax: GitLab only supports * and ** wildcards",
				IncludePath: localInclude.Path,
			})
			continue
		}

		includePaths, err := localInclude.ResolvePaths(projectRoot, entryFilePath)
		if err != nil {
			return line, err
		}

		checkAndWarnEmptyGlobPattern(virtualFile, localInclude, includePaths)

		line = addIncludeParts(virtualFile, includePaths, &addedIncludePaths, combined, line)
	}

	return line, nil
}

// checkAndWarnEmptyGlobPattern adds a finding if a glob pattern resolves to no files
func checkAndWarnEmptyGlobPattern(virtualFile *VirtualCiYamlFile, localInclude *includes2.LocalInclude, includePaths []string) {
	if localInclude.IsGlobPattern() && len(includePaths) == 0 {
		virtualFile.ResolveFindings = append(virtualFile.ResolveFindings, VirtualFileResolveFinding{
			Code:        101,
			Severity:    1, // Warning level
			Message:     "Include pattern '" + localInclude.Path + "' did not match any files",
			IncludePath: localInclude.Path,
		})
	}
}

// addIncludeParts adds parts for all resolved include paths
func addIncludeParts(virtualFile *VirtualCiYamlFile, includePaths []string, addedIncludePaths *[]string, combined *bytes.Buffer, startLine int) int {
	line := startLine

	for _, includePath := range includePaths {
		if slices.Contains(*addedIncludePaths, includePath) {
			continue
		}

		part, err := createPart(includePath, line)
		if err != nil {
			// Add error as a finding to be reported later
			virtualFile.ResolveFindings = append(virtualFile.ResolveFindings, VirtualFileResolveFinding{
				Code:        102,
				Severity:    0, // Error level
				Message:     "Include file '" + includePath + "' could not be loaded: " + err.Error(),
				IncludePath: includePath,
			})
			continue
		}

		*addedIncludePaths = append(*addedIncludePaths, includePath)
		virtualFile.Parts = append(virtualFile.Parts, part)
		combined.Write(part.CiYamlFile.FileContent)
		combined.Write([]byte("\n"))
		line = part.EndLine + 1
	}

	return line
}

// finalizeCombinedYaml adds the include block and creates the combined YAML file
func finalizeCombinedYaml(virtualFile *VirtualCiYamlFile, uniqueLocalIncludes []includes2.Include, combined *bytes.Buffer) error {
	combined.Write([]byte("\n# Auto generated include block\n"))
	includeBlockYaml, err := joinIncludes(uniqueLocalIncludes)
	if err != nil {
		return err
	}
	combined.Write(includeBlockYaml)

	combinedCiYaml, err := NewCiYamlFile(combined.Bytes())
	if err != nil {
		return err
	}
	virtualFile.Combined = combinedCiYaml

	return nil
}

func createPart(includePath string, line int) (*VirtualCiYamlFilePart, error) {
	includeContent, err := os.ReadFile(includePath)
	if err != nil {
		return nil, err
	}

	includeCiYaml, err := NewCiYamlFile(includeContent)
	if err != nil {
		return nil, err
	}

	includeLineCount := bytes.Count(includeContent, []byte("\n"))
	part := &VirtualCiYamlFilePart{
		CiYamlFile: includeCiYaml,
		StartLine:  line,
		EndLine:    line + includeLineCount,
		Path:       includePath,
	}
	return part, nil
}

func joinIncludes(uniqueLocalIncludes []includes2.Include) ([]byte, error) {
	includeNodes := make([]*yaml.Node, len(uniqueLocalIncludes))
	for idx, i := range uniqueLocalIncludes {
		includeNodes[idx] = i.Node()
	}
	includeBlockYaml, err := yaml.Marshal(map[string][]*yaml.Node{
		"include": includeNodes,
	})
	if err != nil {
		return nil, err
	}
	return includeBlockYaml, err
}
