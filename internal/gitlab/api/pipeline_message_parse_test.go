package api

import "testing"

func TestParsePipelineMessage(t *testing.T) {
	testCases := []struct {
		name                 string
		message              string
		expectedJob          string
		expectedErrorMessage string
	}{
		{
			name:                 "generic error message for jobs",
			message:              "jobs deploy config should implement the script:, run:, or trigger: keyword",
			expectedJob:          "jobs",
			expectedErrorMessage: "deploy config should implement the script:, run:, or trigger: keyword",
		},
		{
			name:                 "invalid job",
			message:              "jobs:deploy config contains unknown keys: varia1blaes",
			expectedJob:          "deploy",
			expectedErrorMessage: "config contains unknown keys: varia1blaes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job, errorMessage := ParsePipelineMessage(tc.message)

			if job != tc.expectedJob {
				t.Fatalf("Expected job %s, but got %s", tc.expectedJob, job)
			}

			if errorMessage != tc.expectedErrorMessage {
				t.Fatalf("Expected error message %s, but got %s", tc.expectedErrorMessage, errorMessage)
			}
		})
	}
}
