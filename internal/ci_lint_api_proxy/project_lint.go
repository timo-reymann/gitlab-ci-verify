package ci_lint_api_proxy

import (
	"io"
	"net/http"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/logging"
)

func projectLintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErr(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}

	project := r.PathValue("project")
	ctx := getApiContext(r)
	body, err := io.ReadAll(r.Body)

	if err != nil {
		logging.Warn("Failed to read request body for project", project, err)
		sendErr(w, http.StatusInternalServerError, "Failed to read request body")
		return
	}

	req, err := ctx.GitLabClient.NewRequestWithContext(r.Context(), "POST", "/api/v4/projects/"+project+"/ci/lint", body)
	if err != nil {
		logging.Warn("Failed to create GitLab request for project", project, err)
		sendErr(w, http.StatusInternalServerError, "Failed to prepare GitLab request")
		return
	}

	result, err := ctx.GitLabClient.Do(req)
	if err != nil {
		logging.Warn("Failed to execute GitLab request for project", project, "with URL", req.URL.String(), err)
		sendErr(w, http.StatusInternalServerError, "Failed to request GitLab API: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = io.Copy(w, result.Body)
}
