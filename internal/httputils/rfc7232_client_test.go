package httputils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
			reader, err := client.ReadFileWithCondition(ts.URL, tc.etag, tc.lastModified)
			if (err != nil) != tc.wantErr {
				t.Errorf("DownloadFileWithCondition() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.expectDownload {
				if _, err := io.ReadAll(reader); os.IsNotExist(err) {
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
