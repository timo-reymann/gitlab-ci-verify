package formatter

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/checks"
	"testing"
)

func TestGitLabCodeQualityFormatter(t *testing.T) {
	// Create a test finding
	finding := &checks.CheckFinding{
		Severity: checks.SeverityError,
		Code:     "TEST-001",
		Line:     10,
		Message:  "Test finding message",
		Link:     "https://example.com/test",
		File:     "/path/to/file.yml",
	}

	// Create formatter and buffer
	var buf bytes.Buffer
	formatter := &GitLabCodeQualityFormatter{}

	// Initialize formatter
	err := formatter.Init(&buf)
	assert.NoError(t, err)

	// Start formatting
	err = formatter.Start()
	assert.NoError(t, err)

	// Print finding
	err = formatter.Print(finding)
	assert.NoError(t, err)

	// End formatting
	err = formatter.End()
	assert.NoError(t, err)

	// Parse the output
	var results []map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &results)
	assert.NoError(t, err)

	// Verify the structure
	assert.Len(t, results, 1)
	result := results[0]

	assert.Equal(t, "Test finding message", result["description"])
	assert.Equal(t, "critical", result["severity"])
	assert.NotEmpty(t, result["fingerprint"])

	location := result["location"].(map[string]interface{})
	assert.Equal(t, "/path/to/file.yml", location["path"])

	lines := location["lines"].(map[string]interface{})
	assert.Equal(t, float64(10), lines["begin"])
}

func TestGitLabCodeQualityFormatterMultipleFindings(t *testing.T) {
	// Create test findings
	findings := []*checks.CheckFinding{
		{
			Severity: checks.SeverityError,
			Code:     "TEST-001",
			Line:     5,
			Message:  "Critical finding",
			Link:     "https://example.com/critical",
			File:     "/path/to/file1.yml",
		},
		{
			Severity: checks.SeverityInfo,
			Code:     "TEST-002",
			Line:     15,
			Message:  "Low severity finding",
			Link:     "https://example.com/low",
			File:     "/path/to/file2.yml",
		},
	}

	// Create formatter and buffer
	var buf bytes.Buffer
	formatter := &GitLabCodeQualityFormatter{}

	// Initialize formatter
	err := formatter.Init(&buf)
	assert.NoError(t, err)

	// Start formatting
	err = formatter.Start()
	assert.NoError(t, err)

	// Print findings
	for _, finding := range findings {
		err = formatter.Print(finding)
		assert.NoError(t, err)
	}

	// End formatting
	err = formatter.End()
	assert.NoError(t, err)

	// Parse the output
	var results []map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &results)
	assert.NoError(t, err)

	// Verify the structure
	assert.Len(t, results, 2)

	// Check first finding
	result1 := results[0]
	assert.Equal(t, "Critical finding", result1["description"])
	assert.Equal(t, "critical", result1["severity"])

	// Check second finding
	result2 := results[1]
	assert.Equal(t, "Low severity finding", result2["description"])
	assert.Equal(t, "info", result2["severity"])
}

func TestGitLabCodeQualityFormatterGet(t *testing.T) {
	formatter, err := Get("gitlab")
	assert.NoError(t, err)
	assert.IsType(t, &GitLabCodeQualityFormatter{}, formatter)
}
