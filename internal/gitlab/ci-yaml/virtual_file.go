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
	// Warnings contains warnings generated during virtual file creation (e.g., empty glob patterns)
	Warnings []VirtualFileWarning
}

// VirtualFileWarning represents a warning generated during virtual file creation
type VirtualFileWarning struct {
	// Message is the warning message
	Message string
	// IncludePath is the path of the include that generated the warning
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
		EntryFile:     entryFile,
		EntryFilePath: entryFilePath,
		Parts:         []*VirtualCiYamlFilePart{},
		Warnings:      []VirtualFileWarning{},
	}

	localIncludes := includes2.FilterByTypes(virtualFile.EntryFile.Includes, "local")
	uniqueLocalIncludes := includes2.Unique(localIncludes)

	line := 0
	entryLineCount := bytes.Count(entryFile.FileContent, []byte("\n"))
	combined := bytes.NewBuffer(entryFile.FileContent)
	virtualFile.Parts = append(virtualFile.Parts, &VirtualCiYamlFilePart{
		CiYamlFile: entryFile,
		Path:       entryFilePath,
		StartLine:  1,
		EndLine:    entryLineCount + 2,
	})
	combined.Write([]byte("\n"))
	line += entryLineCount + 2

	addedIncludePaths := make([]string, 0)
	for _, uniqueLocalInclude := range uniqueLocalIncludes {
		localInclude := uniqueLocalInclude.(*includes2.LocalInclude)
		includePaths, err := localInclude.ResolvePaths(projectRoot, entryFilePath)
		if err != nil {
			return nil, err
		}

		// Warn if glob pattern resolves to no files
		if localInclude.IsGlobPattern() && len(includePaths) == 0 {
			virtualFile.Warnings = append(virtualFile.Warnings, VirtualFileWarning{
				Message:     "Glob pattern did not match any files",
				IncludePath: localInclude.Path,
			})
		}

		for _, includePath := range includePaths {
			if slices.Contains(addedIncludePaths, includePath) {
				continue
			}

			part, err := createPart(includePath, line)
			if err != nil {
				return nil, err
			}

			addedIncludePaths = append(addedIncludePaths, includePath)
			virtualFile.Parts = append(virtualFile.Parts, part)
			combined.Write(part.CiYamlFile.FileContent)
			combined.Write([]byte("\n"))
			line = part.EndLine + 1
		}
	}

	combined.Write([]byte("\n# Auto generated include block\n"))
	includeBlockYaml, err := joinIncludes(uniqueLocalIncludes)
	if err != nil {
		return nil, err
	}
	combined.Write(includeBlockYaml)

	combinedCiYaml, err := NewCiYamlFile(combined.Bytes())
	if err != nil {
		return nil, err
	}
	virtualFile.Combined = combinedCiYaml

	return virtualFile, nil
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
