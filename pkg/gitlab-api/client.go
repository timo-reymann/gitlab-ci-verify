package gitlab_api

import (
	"bytes"
	"net/http"
	"strings"
	"time"
)

type GitLabApiClient struct {
	apiBaseUrl string
	token      string
	httpClient http.Client
}

type GitlabApiRequest struct {
	path string
	*http.Request
}

func (g *GitLabApiClient) NewRequest(method string, path string, payload []byte) (*GitlabApiRequest, error) {
	strippedPath := strings.TrimPrefix(path, "/")
	req, err := http.NewRequest(method, g.apiBaseUrl+"/"+strippedPath, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	return &GitlabApiRequest{
		path:    strippedPath,
		Request: req,
	}, nil
}

func (g *GitLabApiClient) Do(r *GitlabApiRequest) (*http.Response, error) {
	return g.httpClient.Do(r.Request)
}

func NewGitlabApiClient(baseUrl string, token string) *GitLabApiClient {
	return &GitLabApiClient{
		apiBaseUrl: strings.TrimSuffix(baseUrl, "/") + "/",
		token:      token,
		httpClient: http.Client{
			Transport: nil,
			Timeout:   5 * time.Second,
		},
	}
}
