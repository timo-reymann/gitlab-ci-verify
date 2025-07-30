package httputils

import (
	"bytes"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/cache"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/hashing"
	"io"
	"net/http"
	"strings"
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

func (hc *RFC7232HttpClient) getCachedMeta(urlHash string) (string, *time.Time, error) {
	if !cache.Exists(urlHash+".meta") || !cache.Exists(urlHash) {
		return "", nil, nil
	}

	metaReader, err := cache.OpenFile(urlHash + ".meta")
	if err != nil {
		return "", nil, err
	}

	meta, err := io.ReadAll(metaReader)
	if err != nil {
		return "", nil, err
	}

	content := strings.Split(string(meta), ";")
	etag := content[0]
	lastModified, err := time.Parse(http.TimeFormat, content[1])
	if err != nil {
		return "", nil, err
	}
	return etag, &lastModified, nil
}

func (hc *RFC7232HttpClient) ReadRemoteOrCached(url string) (io.ReadCloser, error) {
	urlHash := hashing.CreateHashFromString(url)

	etag, lastModified, err := hc.getCachedMeta(urlHash)
	if err != nil {
		return nil, err
	}

	if lastModified == nil {
		lastModified = &time.Time{}
	}

	res, err := hc.GetWithCondition(url, etag, *lastModified)
	if err != nil {
		return nil, err
	} else if res != nil {
		_ = hc.setCachedMeta(res, urlHash)
		return res.Body, nil
	}

	return cache.OpenFile(urlHash)
}

func (hc *RFC7232HttpClient) setCachedMeta(res *http.Response, urlHash string) error {
	metaBuffer := bytes.NewBuffer([]byte(res.Header.Get("ETag") + ";" + res.Header.Get("Last-Modified")))
	return cache.WriteFile(urlHash+".meta", metaBuffer)
}

// GetWithCondition reads a file from an HTTP endpoint using conditional headers
func (hc *RFC7232HttpClient) GetWithCondition(url string, etag string, lastModified time.Time) (*http.Response, error) {
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

	return res, nil
}
