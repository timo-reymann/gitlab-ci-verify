package httputils

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// RFC7232HttpClient is a custom client with support for RFC7232 conditional headers
type RFC7232HttpClient struct {
	client *http.Client
}

// NewRfc7232HttpClient initializes a new RFC7232HttpClient using
// the NewRetryableClient function
func NewRfc7232HttpClient() *RFC7232HttpClient {
	return &RFC7232HttpClient{
		client: NewRetryableClient(),
	}
}

// ReadFileWithCondition reads a file from an HTTP endpoint using conditional headers
func (hc *RFC7232HttpClient) ReadFileWithCondition(url, etag string, lastModified time.Time) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}

	if !lastModified.IsZero() {
		req.Header.Set("If-Modified-Since", lastModified.Format(http.TimeFormat))
	}

	res, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if res.StatusCode == http.StatusNotModified {
		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res.Body, nil
}
