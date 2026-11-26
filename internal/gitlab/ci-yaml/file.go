package ci_yaml

import (
	"bytes"
	"regexp"

	formatconversion "github.com/timo-reymann/gitlab-ci-verify/v2/internal/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/ci-yaml/includes"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/filtering"
	"gopkg.in/yaml.v3"
)

// CiYamlFile represents a parsed CI YAML file
type CiYamlFile struct {
	// FileContent contains the raw file content
	FileContent []byte
	// ParsedYamlMap contains the parsed YAML configuration (without the header) as map structure
	ParsedYamlMap map[string]any
	// ParsedYamlDoc contains the parsed YAML configuration (without the header) as node structure
	ParsedYamlDoc *yaml.Node
	// ParsedHeaderDoc contains the parsed header document (`spec:`) as node structure, nil if the file has no header
	ParsedHeaderDoc *yaml.Node
	// Includes contains all includes in the entry file
	Includes []includes.Include
	// LineNumberMapping contains the line number node mapping for the parsed YAML
	lineNumberMapping yaml.LineNumberMapping
}

// HasHeader checks if the file has a header document, e.g., the `spec:` section of a CI/CD component
func (cyf *CiYamlFile) HasHeader() bool {
	return cyf.ParsedHeaderDoc != nil
}

func firstKeyHeadComment(doc *yaml.Node) string {
	if doc == nil || len(doc.Content) == 0 || len(doc.Content[0].Content) == 0 {
		return ""
	}
	return doc.Content[0].Content[0].HeadComment
}

func (cyf *CiYamlFile) GetFileLevelIgnores() []filtering.IgnoreComment {
	ignores := filtering.ParseIgnoreComment(firstKeyHeadComment(cyf.ParsedHeaderDoc))
	return append(ignores, filtering.ParseIgnoreComment(firstKeyHeadComment(cyf.ParsedYamlDoc))...)
}

func (cyf *CiYamlFile) GetLineLevelIgnores(line int) []filtering.IgnoreComment {
	return filtering.ParseIgnoreForLine(cyf.lineNumberMapping, []int{line})
}

// hasSpecKey checks if the document is a mapping containing the `spec` keyword
func hasSpecKey(doc *yaml.Node) bool {
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return false
	}

	root := doc.Content[0]
	for idx := 0; idx < len(root.Content); idx += 2 {
		if root.Content[idx].Value == "spec" {
			return true
		}
	}

	return false
}

// selectDocuments follows the GitLab semantics for multi document CI YAML files.
//
// Empty documents are ignored. In case there are at least two documents and the first one contains the `spec`
// keyword, the first one is the header and the second one the configuration. Otherwise, the first document is the
// configuration.
//
// It returns the index of the header document (-1 if there is none) and the index of the configuration document.
func selectDocuments(docs []*yaml.Node) (headerIdx int, contentIdx int) {
	nonEmpty := make([]int, 0, len(docs))
	for idx, doc := range docs {
		if !formatconversion.IsEmptyYamlDocument(doc) {
			nonEmpty = append(nonEmpty, idx)
		}
	}

	switch {
	case len(nonEmpty) == 0:
		return -1, 0
	case len(nonEmpty) >= 2 && hasSpecKey(docs[nonEmpty[0]]):
		return nonEmpty[0], nonEmpty[1]
	default:
		return -1, nonEmpty[0]
	}
}

// decodeDocumentToMap decodes the document with the given index of the stream to a map
func decodeDocumentToMap(content []byte, docIdx int) (map[string]any, error) {
	decoder := yaml.NewDecoder(bytes.NewBuffer(content))
	decoder.UniqueKeys(false)

	for idx := 0; idx < docIdx; idx++ {
		var skipped any
		if err := decoder.Decode(&skipped); err != nil {
			return nil, err
		}
	}

	var parsed map[string]any
	if err := decoder.Decode(&parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}

// NewCiYamlFile from byte contents
//
// Files containing a header (`spec:`) followed by the configuration in a second document are supported, see
// selectDocuments for details.
func NewCiYamlFile(content []byte) (*CiYamlFile, error) {
	docs, lineNumberMapping, err := formatconversion.ParseYamlDocuments(content)
	if err != nil {
		return nil, err
	}

	headerIdx, contentIdx := selectDocuments(docs)
	doc := docs[contentIdx]

	var header *yaml.Node
	if headerIdx != -1 {
		header = docs[headerIdx]
	}

	parsed, err := decodeDocumentToMap(content, contentIdx)
	if err != nil {
		return nil, err
	}

	parsedIncludes, err := includes.ParseIncludes(doc)
	if err != nil {
		return nil, err
	}

	return &CiYamlFile{
		FileContent:       content,
		ParsedYamlMap:     parsed,
		ParsedYamlDoc:     doc,
		ParsedHeaderDoc:   header,
		Includes:          parsedIncludes,
		lineNumberMapping: lineNumberMapping,
	}, nil
}

var documentMarker = regexp.MustCompile(`^---(\s|$)`)

// ConfigurationContent returns the file content with everything that is not part of the configuration document
// replaced by empty lines. This includes the header document (`spec:`), document markers and any following
// documents. The line count is preserved, so line numbers stay the same as in FileContent.
//
// This allows to concat the configuration of multiple files into one YAML document.
func (cyf *CiYamlFile) ConfigurationContent() []byte {
	lines := bytes.SplitAfter(cyf.FileContent, []byte("\n"))

	startLine := 0 // 1-based line where the configuration document starts, 0 = no content
	if cyf.ParsedYamlDoc != nil && len(cyf.ParsedYamlDoc.Content) > 0 {
		startLine = cyf.ParsedYamlDoc.Content[0].Line
	}

	// blank everything up to and including the last document marker before the configuration
	blankUntil := 0
	for idx := 0; idx < len(lines) && (startLine == 0 || idx+1 < startLine); idx++ {
		if documentMarker.Match(lines[idx]) {
			blankUntil = idx + 1
		}
	}

	// blank everything starting with the first document marker after the configuration
	blankFrom := len(lines)
	for idx := max(startLine, blankUntil); idx < len(lines); idx++ {
		if documentMarker.Match(lines[idx]) {
			blankFrom = idx
			break
		}
	}

	result := bytes.NewBuffer(make([]byte, 0, len(cyf.FileContent)))
	for idx, line := range lines {
		if idx < blankUntil || idx >= blankFrom {
			if bytes.HasSuffix(line, []byte("\n")) {
				result.WriteByte('\n')
			}
			continue
		}
		result.Write(line)
	}

	return result.Bytes()
}
