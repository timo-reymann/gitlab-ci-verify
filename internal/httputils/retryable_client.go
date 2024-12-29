package httputils

import (
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"time"
)

// NewRetryableClient creates a new http client with
// retryablehttp.Client as transport
func NewRetryableClient() *http.Client {
	retryHttpClient := retryablehttp.NewClient()
	retryHttpClient.RetryMax = 3
	retryHttpClient.HTTPClient.Timeout = 5 * time.Second
	return retryHttpClient.StandardClient()
}
