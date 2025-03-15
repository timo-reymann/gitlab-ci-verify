package httputils

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"net/http"
	"time"
)

type httpLogger struct{}

func (h httpLogger) Error(msg string, keysAndValues ...interface{}) {
	logging.Debug(msg, keysAndValues)
}

func (h httpLogger) Info(msg string, keysAndValues ...interface{}) {
	logging.Verbose(msg, keysAndValues)
}

func (h httpLogger) Debug(msg string, keysAndValues ...interface{}) {
	logging.Debug(msg, keysAndValues)
}

func (h httpLogger) Warn(msg string, keysAndValues ...interface{}) {
	logging.Warn(msg, keysAndValues)
}

// NewRetryableClient creates a new http client with
// retryablehttp.Client as transport
func NewRetryableClient() *http.Client {
	retryHttpClient := retryablehttp.NewClient()
	retryHttpClient.RetryMax = 3
	retryHttpClient.HTTPClient.Timeout = 5 * time.Second
	retryHttpClient.Logger = &httpLogger{}
	return retryHttpClient.StandardClient()
}
