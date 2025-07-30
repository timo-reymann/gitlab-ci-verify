package httputils

import (
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/cache"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/hashing"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func mockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		etag := "\"123456789abcdef\""
		lastModified := "Fri, 01 Dec 2023 00:00:00 GMT"

		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		if r.Header.Get("If-Modified-Since") == lastModified {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Header().Set("ETag", etag)
		w.Header().Set("Last-Modified", lastModified)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("file content"))
	}))
}

func TestRFC7232HttpClient_DownloadFileWithCondition(t *testing.T) {
	ts := mockServer()
	defer ts.Close()

	client := NewRfc7232HttpClient()

	testCases := []struct {
		name           string
		etag           string
		lastModified   time.Time
		expectDownload bool
		wantErr        bool
	}{
		{
			name:           "File should be downloaded",
			etag:           "",
			lastModified:   time.Time{},
			expectDownload: true,
			wantErr:        false,
		},
		{
			name:           "File should not be downloaded (HTTP 304)",
			etag:           "\"123456789abcdef\"",
			lastModified:   time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			expectDownload: false,
			wantErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader, err := client.GetWithCondition(ts.URL, tc.etag, tc.lastModified)
			if (err != nil) != tc.wantErr {
				t.Errorf("DownloadFileWithCondition() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.expectDownload {
				if _, err := io.ReadAll(reader.Body); os.IsNotExist(err) {
					t.Errorf("Expected file to exist, but it does not")
				}
			} else {
				if reader != nil {
					t.Errorf("Expected file not to be downloaded, but it exists")
				}
			}
		})
	}
}

func TestReadRemoteOrCached(t *testing.T) {
	ts := mockServer()
	defer ts.Close()

	urlhash := hashing.CreateHashFromString(ts.URL)

	client := NewRfc7232HttpClient()

	testCases := []struct {
		name           string
		setupCache     func()
		expectDownload bool
		expectErr      bool
	}{
		{
			name: "Cache contains meta, but no file",
			setupCache: func() {
				cache.WriteFile(urlhash+".meta", strings.NewReader("\"123456789abcdef\";Fri, 01 Dec 2023 00:00:00 GMT"))
			},
			expectDownload: true,
			expectErr:      false,
		},
		{
			name: "Cache contains meta and file",
			setupCache: func() {
				println(urlhash)
				cache.WriteFile(urlhash+".meta", strings.NewReader("\"123456789abcdef\";Fri, 01 Dec 2023 00:00:00 GMT"))
				cache.WriteFile(urlhash, strings.NewReader("cached file content"))
			},
			expectDownload: false,
			expectErr:      false,
		},
		{
			name: "Cache contains none and downloads",
			setupCache: func() {

			},
			expectDownload: true,
			expectErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean cache
			dir, _ := cache.CacheDir()
			os.RemoveAll(dir)
			tc.setupCache()

			reader, err := client.ReadRemoteOrCached(ts.URL)
			if (err != nil) != tc.expectErr {
				t.Errorf("ReadRemoteOrCached() error = %v, expectErr %v", err, tc.expectErr)
			}

			if tc.expectDownload {
				if reader == nil {
					t.Errorf("Expected file to be downloaded, but it was not")
				} else {
					content, _ := io.ReadAll(reader)
					if string(content) != "file content" {
						t.Errorf("Expected downloaded content to be 'file content', got %s", string(content))
					}
				}
			} else {
				if reader == nil {
					t.Errorf("Expected file to be read from cache, but it was not")
				} else {
					content, _ := io.ReadAll(reader)
					if string(content) != "cached file content" {
						t.Errorf("Expected cached content to be 'cached file content', got '%s'", string(content))
					}
				}
			}
		})
	}
}
