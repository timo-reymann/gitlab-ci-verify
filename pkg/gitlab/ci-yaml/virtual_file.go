package ci_yaml

import (
	"bytes"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml/includes"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/location"
	"gopkg.in/yaml.v3"
	"os"
	"slices"
	"sort"
)

type VirtualCiYamlFilePart struct {
	CiYamlFile *CiYamlFile
	Path       string
	StartLine  int
	EndLine    int
}

type VirtualCiYamlFile struct {
	EntryFile     *CiYamlFile
	EntryFilePath string
	Combined      *CiYamlFile
	Parts         []*VirtualCiYamlFilePart
}

func (v *VirtualCiYamlFile) Resolve(line int) (*VirtualCiYamlFilePart, *location.Location) {
	idx := sort.Search(len(v.Parts), func(i int) bool {
		return v.Parts[i].StartLine > line
	})
	var part *VirtualCiYamlFilePart
	if idx > 0 {
		part = v.Parts[idx-1]
	}

	if part != nil && part.StartLine <= line && part.EndLine >= line {
		return v.Parts[idx-1], location.NewLocation(part.Path, line-part.StartLine)
	}

	return nil, nil
}

func CreateVirtualCiYamlFile(projectRoot string, entryFilePath string, entryFile *CiYamlFile) (*VirtualCiYamlFile, error) {
	virtualFile := &VirtualCiYamlFile{
		EntryFile:     entryFile,
		EntryFilePath: entryFilePath,
		Parts:         []*VirtualCiYamlFilePart{},
	}

	localIncludes := includes.FilterByTypes(virtualFile.EntryFile.Includes, "local")
	uniqueLocalIncludes := includes.Unique(localIncludes)

	line := 0
	entryLineCount := bytes.Count(entryFile.FileContent, []byte("\n"))
	combined := bytes.NewBuffer(entryFile.FileContent)
	line += entryLineCount + 1

	addedIncludePaths := make([]string, 0)
	for _, uniqueLocalInclude := range uniqueLocalIncludes {
		localInclude := uniqueLocalInclude.(*includes.LocalInclude)
		includePath := localInclude.ResolvePath(projectRoot, entryFilePath)

		if slices.Contains(addedIncludePaths, includePath) {
			continue
		}

		part, err := createPart(includePath, line)
		if err != nil {
			return nil, err
		}

		virtualFile.Parts = append(virtualFile.Parts, part)
		combined.Write(part.CiYamlFile.FileContent)
		line += part.EndLine + 1
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

func joinIncludes(uniqueLocalIncludes []includes.Include) ([]byte, error) {
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
