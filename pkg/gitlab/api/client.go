package api

import (
	"context"
	"fmt"
	format_conversion "github.com/timo-reymann/gitlab-ci-verify/pkg/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/logging"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/netrc"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Client to access the gitlab api
type Client struct {
	apiBaseUrl string
	token      string
	httpClient http.Client
}

// Do the http request to the gitlab api
func (g *Client) Do(r *Request) (*Response, error) {
	logging.Debug("perform gitlab api request", r.Method, r.path, "against API host", g.apiBaseUrl)
	res, err := g.httpClient.Do(r.Request)
	if err != nil {
		return nil, err
	}

	return &Response{res}, err
}

// NewRequest creates a new request prefixing it with the base url and adding the GitLab access token
func (g *Client) NewRequest(method string, path string, payload []byte) (*Request, error) {
	return NewRequest(context.Background(), method, g.apiBaseUrl, path, g.token, payload)
}

// NewRequestWithContext creates a new request prefixing it with the base url and adding the GitLab access token bound to a context
func (g *Client) NewRequestWithContext(ctx context.Context, method string, path string, payload []byte) (*Request, error) {
	return NewRequest(ctx, method, g.apiBaseUrl, path, g.token, payload)
}

// LintCiYaml against the api for a project
func (g *Client) LintCiYaml(ctx context.Context, projectSlug string, ciYaml []byte) (*CiLintResult, error) {
	wrapped, err := format_conversion.ToJson(map[string]any{
		"content": string(ciYaml),
	})
	if err != nil {
		return nil, err
	}

	req, err := g.NewRequestWithContext(ctx, "POST", fmt.Sprintf("/api/v4/projects/%s/ci/lint", url.QueryEscape(projectSlug)), wrapped)
	if err != nil {
		return nil, err
	}

	res, err := g.Do(req)
	if err != nil || res.CheckStatus() != nil {
		return nil, err
	}

	var lintResult = &CiLintResult{}
	if err = res.UnmarshalJson(lintResult); err != nil {
		return nil, err
	}

	return lintResult, nil
}

// NewClient creates a new api client instance for the gitlab api
func NewClient(baseUrl string, token string) *Client {
	protocolPrefix := ""
	if !strings.HasPrefix(baseUrl, "http") {
		protocolPrefix = "https://"
	}

	return &Client{
		apiBaseUrl: protocolPrefix + strings.TrimSuffix(baseUrl, "/") + "/",
		token:      token,
		httpClient: http.Client{
			Transport: nil,
			Timeout:   5 * time.Second,
		},
	}
}

// NewClientWithMultiTokenSources creates a new client instance for the gitlab api, taking multiple sources for the token
// The order is as following
// 1. token specified via parameter is non-empty
// 2. netrc contains a host entry with a password
// 3. environment variable GITLAB_TOKEN is set
func NewClientWithMultiTokenSources(baseUrl string, token string) *Client {
	if token != "" {
		return NewClient(baseUrl, token)
	}

	// try to load netrc for user and get credential form that
	userNetrc, err := netrc.ReadUserNetrc()
	if err == nil {
		credentials, err := netrc.GetCredentials(userNetrc, baseUrl)
		if err == nil && credentials.Password != "" {
			return NewClient(baseUrl, credentials.Password)
		}
	}

	// fallback is the GITLAB_TOKEN env var
	return NewClient(baseUrl, os.Getenv("GITLAB_TOKEN"))
}
