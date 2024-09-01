package ci_yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestExtractScripts(t *testing.T) {
	testCases := []struct {
		gitlabCiYamlName string
		scripts          map[string]struct {
			scriptNodeCount map[string]int
		}
	}{
		{
			gitlabCiYamlName: "validSingleJobSingleLine",
			scripts: map[string]struct{ scriptNodeCount map[string]int }{
				"build-job": {
					scriptNodeCount: map[string]int{
						"script": 1,
					},
				},
			},
		},
		{
			gitlabCiYamlName: "validSingleJobMultiLine",
			scripts: map[string]struct{ scriptNodeCount map[string]int }{
				"build-job": {
					scriptNodeCount: map[string]int{
						"script": 2,
					},
				},
			},
		},
		{
			gitlabCiYamlName: "validMultiJob",
			scripts: map[string]struct{ scriptNodeCount map[string]int }{
				"build-job": {
					scriptNodeCount: map[string]int{
						"script": 2,
					},
				},
				"test-job": {
					scriptNodeCount: map[string]int{
						"script": 1,
					},
				},
				"test-job-with-before-script": {
					map[string]int{
						"script":        1,
						"before_script": 1,
					},
				},
			},
		},
		{
			gitlabCiYamlName: "validSingleJobMultipleScripts",
			scripts: map[string]struct{ scriptNodeCount map[string]int }{
				"job-name": {
					map[string]int{
						"script":        1,
						"before_script": 1,
						"after_script":  1,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.gitlabCiYamlName, func(t *testing.T) {
			content, err := os.ReadFile(fmt.Sprintf("test_data/gitlab-ci/%s.yaml", tc.gitlabCiYamlName))
			if err != nil {
				t.Fatal(err)
			}
			var node yaml.Node
			if err = yaml.Unmarshal(content, &node); err != nil {
				t.Fatal(err.Error())
			}

			resultCount := 0

			for script := range ExtractScripts(&node) {
				resultCount++
				if script.JobName == "" {
					t.Fatalf("Expected non-empty job name")
				}

				expectedScripts, ok := tc.scripts[script.JobName]
				if !ok {
					t.Fatalf("Unexpected job %s parsed", script.JobName)
				}

				for key, expectedNodeCount := range expectedScripts.scriptNodeCount {
					if expectedNodeCount != len(script.ScriptParts[key]) {
						t.Fatalf("Expected %d script parts got %d", expectedNodeCount, len(script.ScriptParts[key]))
					}
				}

				fmt.Printf("%v", script)
			}

			if resultCount != len(tc.scripts) {
				t.Fatalf("Expected %d results, got %d", len(tc.scripts), resultCount)
			}
		})
	}
}
