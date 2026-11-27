package format_conversion

import (
	"bytes"
	"errors"
	"io"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/logging"
	"gopkg.in/yaml.v3"
)

// ParseYamlNode from the given raw content
//
// Only the first document of a multi document stream is returned, use ParseYamlDocuments to get all documents.
func ParseYamlNode(content []byte) (*yaml.Node, error) {
	var result yaml.Node

	decoder := yaml.NewDecoder(bytes.NewBuffer(content))
	decoder.UniqueKeys(false)
	err := decoder.Decode(&result)
	return &result, err
}

// ParseYamlDocuments parses all documents of a (multi document) YAML stream.
//
// When the stream contains no document at all, io.EOF is returned, mirroring the behavior of ParseYamlNode.
// The returned line number mapping covers all documents, line numbers are relative to the start of the stream.
func ParseYamlDocuments(content []byte) ([]*yaml.Node, yaml.LineNumberMapping, error) {
	decoder := yaml.NewDecoder(bytes.NewBuffer(content))
	decoder.UniqueKeys(false)
	decoder.WithLineNumberMapping()

	docs := make([]*yaml.Node, 0, 1)
	for {
		var doc yaml.Node
		err := decoder.Decode(&doc)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			// Only the first document is mandatory, failing to parse following documents keeps the documents
			// parsed so far. This is the case, e.g., for flow style (JSON) files that are followed by block style
			// content.
			if len(docs) > 0 {
				logging.Debug("Ignoring YAML content after document", len(docs), "as it failed to parse:", err.Error())
				break
			}
			return nil, nil, err
		}
		docs = append(docs, &doc)
	}

	if len(docs) == 0 {
		return nil, nil, io.EOF
	}

	return docs, decoder.LineNumberMapping(), nil
}

// IsEmptyYamlDocument checks if a document node has no content, e.g., for a stream starting with a document marker
func IsEmptyYamlDocument(doc *yaml.Node) bool {
	if doc == nil || len(doc.Content) == 0 {
		return true
	}

	root := doc.Content[0]
	return root.Kind == yaml.ScalarNode && root.Tag == "!!null"
}
