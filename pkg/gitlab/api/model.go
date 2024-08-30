package api

// CiLintResult represents the response from the Gitlab API when a YAML file has been validated
type CiLintResult struct {
	// Valid indicates if the YAML file is syntactically correct
	Valid bool `json:"valid"`
	// MergedYaml contains the final YAML with all includes and anchors resolved
	MergedYaml string `json:"mergedYaml"`
	// Errors encountered with the pipeline
	Errors []string `json:"errors"`
	// Warnings encountered with the pipeline
	Warnings []string `json:"warnings"`
}
