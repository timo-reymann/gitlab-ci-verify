package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/internal/yamlpathutils"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"strings"
)

type GitlabPagesJobCheck struct{}

func (g GitlabPagesJobCheck) finding(path string, line int) CheckFinding {
	return CheckFinding{
		Severity: SeverityWarning,
		Code:     "GL-201",
		Line:     line,
		Message:  "pages job should contain artifacts with public path",
		Link:     "https://docs.gitlab.com/ee/user/project/pages",
		File:     path,
	}
}

func (g GitlabPagesJobCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	_, jobExists := i.CiYaml.ParsedYamlMap["pages"]
	if !jobExists {
		return []CheckFinding{}, nil
	}

	pagesJob, _ := yamlpathutils.MustPath(yamlpath.NewPath(".pages")).Find(i.CiYaml.ParsedYamlDoc)

	artifactNodes, _ := yamlpathutils.MustPath(yamlpath.NewPath(".artifacts.paths")).Find(pagesJob[0])
	if len(artifactNodes) != 1 {
		return []CheckFinding{
			g.finding(i.Configuration.GitLabCiFile, pagesJob[0].Line-1),
		}, nil
	}

	var paths []string
	if err := artifactNodes[0].Decode(&paths); err != nil {
		return nil, err
	}

	for _, path := range paths {
		if strings.HasPrefix(path, "public") {
			return []CheckFinding{}, nil
		}
	}

	return []CheckFinding{
		g.finding(i.Configuration.GitLabCiFile, artifactNodes[0].Line-1),
	}, nil
}
