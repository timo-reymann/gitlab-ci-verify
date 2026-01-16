package ci_lint_api_proxy

import "net/http"

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	sendErr(w, http.StatusNotFound, "No route found")
}
