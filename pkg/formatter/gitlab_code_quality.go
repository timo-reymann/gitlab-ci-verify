package formatter

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/checks"
)

type GitLabCodeQualityFormatter struct {
	writer           io.Writer
	firstItemWritten bool
	baseDir          string
}

type gitlabCodeQualityFinding struct {
	Description string                    `json:"description"`
	Fingerprint string                    `json:"fingerprint"`
	Severity    string                    `json:"severity"`
	Location    gitlabCodeQualityLocation `json:"location"`
}

type gitlabCodeQualityLocation struct {
	Path  string                 `json:"path"`
	Lines gitlabCodeQualityLines `json:"lines"`
}

type gitlabCodeQualityLines struct {
	Begin int `json:"begin"`
}

func newGitLabCodeQualityFinding(f *checks.CheckFinding, baseDir string) (gitlabCodeQualityFinding, error) {
	loc, err := f.Location()
	if err != nil {
		return gitlabCodeQualityFinding{}, err
	}

	// Make path relative to base directory
	relPath, err := filepath.Rel(baseDir, loc.File)
	if err != nil {
		// If we can't make it relative (e.g., different drives on Windows), use absolute path
		relPath = loc.File
	}

	return gitlabCodeQualityFinding{
		Description: f.Message,
		Fingerprint: f.Fingerprint(),
		Severity:    mapSeverity(f.SeverityName()),
		Location: gitlabCodeQualityLocation{
			Path: relPath,
			Lines: gitlabCodeQualityLines{
				Begin: f.Line,
			},
		},
	}, nil
}

func mapSeverity(severityName string) string {
	switch severityName {
	case "Error":
		return "critical"
	case "Warning":
		return "major"
	case "Info":
		return "info"
	case "Style":
		return "info"
	default:
		return "info"
	}
}

func (g *GitLabCodeQualityFormatter) writeString(val string) error {
	_, err := g.writer.Write([]byte(val))
	return err
}

func (g *GitLabCodeQualityFormatter) Init(w io.Writer) error {
	g.writer = w
	g.firstItemWritten = false

	// Get current working directory to make paths relative
	var err error
	g.baseDir, err = os.Getwd()
	if err != nil {
		return err
	}
	return nil
}

func (g *GitLabCodeQualityFormatter) Start() error {
	return g.writeString("[")
}

func (g *GitLabCodeQualityFormatter) Print(f *checks.CheckFinding) error {
	finding, err := newGitLabCodeQualityFinding(f, g.baseDir)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(finding)
	if err != nil {
		return err
	}

	if g.firstItemWritten {
		if err := g.writeString(","); err != nil {
			return err
		}
	}

	if err := g.writeString("\n"); err != nil {
		return err
	}

	if err := g.writeString("  "); err != nil {
		return err
	}

	_, err = g.writer.Write(buf)
	g.firstItemWritten = true

	return err
}

func (g *GitLabCodeQualityFormatter) End() error {
	return g.writeString("\n]")
}
