package api

import (
	"context"
	"fmt"
	format_conversion "github.com/timo-reymann/gitlab-ci-verify/pkg/format-conversion"
	"net/http"
	"net/url"
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

	req, err := g.NewRequestWithContext(ctx, "POST", fmt.Sprintf("/projects/%s/ci/lint", url.QueryEscape(projectSlug)), wrapped)
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
	if !strings.HasPrefix("http", baseUrl) {
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
