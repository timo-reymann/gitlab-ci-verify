package ci_lint_api_proxy

import "net/http"

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	sendErr(w, http.StatusNotFound, "No route found")
}
