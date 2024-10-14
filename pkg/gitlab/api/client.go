package api

import (
	"context"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"net/http"
	"net/url"
	"reflect"
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

// Prioritized list of token sources
var tokenSources = []TokenSource{
	NetRcTokenSource{},
	VaultTokenSource{},
	EnvVarTokenSource{},
}

// NewClientWithMultiTokenSources creates a new client instance for the gitlab api, taking multiple sources for the token
// The order is as following
// 1. The Token specified via parameter is non-empty
// 2. Try to look up tokenSources and use the first one without an error and return a non-empty string
func NewClientWithMultiTokenSources(baseUrl string, token string) *Client {
	if token != "" && !strings.Contains(token, "://") {
		return NewClient(baseUrl, token)
	}

	hints := TokenSourceLookupHints{
		ExistingToken: token,
		BaseUrl:       baseUrl,
	}

	for _, src := range tokenSources {
		token, err := src.Lookup(hints)
		if err != nil {
			logging.Debug(fmt.Sprintf("Failed to lookup token for token source %v: %s", reflect.TypeOf(tokenSources), err.Error()))
			continue
		}

		if token == "" {
			logging.Debug(fmt.Sprintf("Token for token source %v is empty", reflect.TypeOf(tokenSources)))
			continue
		}

		return NewClient(baseUrl, token)
	}

	// fallback is the GITLAB_TOKEN env var
	return NewClient(baseUrl, token)
}
