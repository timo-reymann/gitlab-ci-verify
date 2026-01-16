package ci_lint_api_proxy

import (
	"net/http"
)

func Serve(endpoint string, ctx ApiContext) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", notFoundHandler)
	mux.HandleFunc("/api/v4/projects/{project}/ci/lint", addContextMiddleware(projectLintHandler, &ctx))
	return http.ListenAndServe(endpoint, mux)
}
