package ci_lint_api_proxy

import (
	"context"
	"net/http"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/api"
)

type ApiContext struct {
	GitLabClient *api.Client
}

const contextKeyApiContext = "apiContext"

func addContextMiddleware(handler http.HandlerFunc, ctx *ApiContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r.WithContext(context.WithValue(r.Context(), contextKeyApiContext, ctx)))
	}
}

func getApiContext(r *http.Request) *ApiContext {
	return r.Context().Value(contextKeyApiContext).(*ApiContext)
}
