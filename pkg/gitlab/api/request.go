package api

import (
	"bytes"
	"context"
	"net/http"
	"strings"
)

// Request wraps a regular http.Request to make it more explicit for the gitlab api client
type Request struct {
	path string
	*http.Request
}

// NewRequest for usage with gitlab api client
func NewRequest(ctx context.Context, method string, baseUrl string, path string, token string, payload []byte) (*Request, error) {
	strippedPath := strings.TrimPrefix(path, "/")
	req, err := http.NewRequestWithContext(ctx, method, baseUrl+"api/v4/"+strippedPath, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")
	return &Request{
		path:    strippedPath,
		Request: req,
	}, nil
}
