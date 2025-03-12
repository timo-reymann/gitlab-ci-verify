package ci_yaml

import (
	"bytes"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml/includes"
	"os"
	"slices"
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

		includeContent, err := os.ReadFile(includePath)
		if err != nil {
			return nil, err
		}

		includeCiYaml, err := NewCiYamlFile(includeContent)
		if err != nil {
			return nil, err
		}

		includeLineCount := bytes.Count(includeContent, []byte("\n"))
		virtualFile.Parts = append(virtualFile.Parts, &VirtualCiYamlFilePart{
			CiYamlFile: includeCiYaml,
			StartLine:  line,
			EndLine:    line + includeLineCount,
			Path:       includePath,
		})
		combined.Write(includeContent)
		line += includeLineCount + 1
	}

	// TODO Add joined includes to end

	combinedCiYaml, err := NewCiYamlFile(combined.Bytes())
	if err != nil {
		return nil, err
	}
	virtualFile.Combined = combinedCiYaml

	return virtualFile, nil
}
