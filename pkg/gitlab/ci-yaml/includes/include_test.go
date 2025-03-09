package includes

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestParseInclude(t *testing.T) {
	testCases := map[string]struct {
		file             string
		expectedIncludes []Include
		expectedErr      error
	}{
		"local qualified": {
			file: "test_data/local_qualified.yml",
			expectedIncludes: []Include{
				NewLocalInclude(nil, "test.yml"),
			},
		},
		"local simple": {
			file: "test_data/local_simple.yml",
			expectedIncludes: []Include{
				NewLocalInclude(nil, "test.yml"),
			},
		},
		"project": {
			file: "test_data/project.yml",
			expectedIncludes: []Include{
				NewProjectInclude(nil, "my-group/my-project", []string{".compliance-gitlab-ci.yml"}),
			},
		},
		"projects": {
			file: "test_data/projects.yml",
			expectedIncludes: []Include{
				NewProjectInclude(nil, "my-group/my-project", []string{"/templates/.gitlab-ci-template.yml"}),
				NewProjectInclude(nil, "my-group/my-subgroup/my-project-2", []string{"/templates/.builds.yml", "/templates/.tests.yml"}),
			},
		},
		"component": {
			file: "test_data/component.yml",
			expectedIncludes: []Include{
				NewComponentInclude(nil, "$CI_SERVER_FQDN/my-org/security-components/secret-detection@1.0"),
			},
		},
		"template": {
			file: "test_data/template.yml",
			expectedIncludes: []Include{
				NewTemplateInclude(nil, "Workflows/MergeRequest-Pipelines.gitlab-ci.yml"),
			},
		},
		"remote": {
			file: "test_data/remote.yml",
			expectedIncludes: []Include{
				NewRemoteInclude(nil, "https://gitlab.com/example-project/-/raw/main/.gitlab-ci.yml", ""),
			},
		},
		"remote with integrity": {
			file: "test_data/remote_with_integrity.yml",
			expectedIncludes: []Include{
				NewRemoteInclude(nil, "https://gitlab.com/example-project/-/raw/main/.gitlab-ci.yml", "sha256-L3/GAoKaw0Arw6hDCKeKQlV1QPEgHYxGBHsH4zG1IY8="),
			},
		},
		"mixed": {
			file: "test_data/mixed.yml",
			expectedIncludes: []Include{
				NewLocalInclude(nil, "test.yml"),
				NewProjectInclude(nil, "my-group/my-project", []string{"/templates/.gitlab-ci-template.yml"}),
				NewProjectInclude(nil, "my-group/my-subgroup/my-project-2", []string{"/templates/.builds.yml", "/templates/.tests.yml"}),
				NewLocalInclude(nil, "test1.yml"),
				NewComponentInclude(nil, "$CI_SERVER_FQDN/my-org/security-components/secret-detection@1.0"),
				NewTemplateInclude(nil, "Workflows/MergeRequest-Pipelines.gitlab-ci.yml"),
				NewRemoteInclude(nil, "https://gitlab.com/example-project/-/raw/main/.gitlab-ci.yml", ""),
			},
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			content, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatal(err)
			}

			var parsedContent yaml.Node
			err = yaml.Unmarshal(content, &parsedContent)
			if err != nil {
				t.Fatal(err)
			}

			include, err := ParseIncludes(&parsedContent)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if len(include) != len(tt.expectedIncludes) {
				t.Fatalf("Expected %d includes, got %d", len(tt.expectedIncludes), len(include))
			}

			if !Equals(include, tt.expectedIncludes) {
				t.Errorf("expected %v, got %v", tt.expectedIncludes, include)
			}
		})
	}
}
