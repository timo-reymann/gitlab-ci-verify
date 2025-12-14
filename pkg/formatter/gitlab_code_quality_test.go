package formatter

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/checks"
	"os"
	"path/filepath"
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
	// The path should now be relative to the current working directory
	// Since /path/to/file.yml doesn't exist, filepath.Rel will create a relative path from current dir
	assert.Contains(t, location["path"], "file.yml")

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

func TestGitLabCodeQualityFormatterRelativePaths(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create a test file in the temp directory
	testFile := filepath.Join(tempDir, "test.yml")
	file, err := os.Create(testFile)
	assert.NoError(t, err)
	file.Close()

	// Change to the temp directory
	originalDir, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	assert.NoError(t, err)

	// Create a test finding with the absolute path
	absPath, err := filepath.Abs(testFile)
	assert.NoError(t, err)

	finding := &checks.CheckFinding{
		Severity: checks.SeverityError,
		Code:     "TEST-001",
		Line:     10,
		Message:  "Test finding message",
		Link:     "https://example.com/test",
		File:     absPath,
	}

	// Create formatter and buffer
	var buf bytes.Buffer
	formatter := &GitLabCodeQualityFormatter{}

	// Initialize formatter (this will get the current working directory)
	err = formatter.Init(&buf)
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
	// The path should now be relative to the current working directory
	assert.Equal(t, "test.yml", location["path"])

	lines := location["lines"].(map[string]interface{})
	assert.Equal(t, float64(10), lines["begin"])
}

func TestGitLabCodeQualityFormatterGet(t *testing.T) {
	formatter, err := Get("gitlab")
	assert.NoError(t, err)
	assert.IsType(t, &GitLabCodeQualityFormatter{}, formatter)
}
