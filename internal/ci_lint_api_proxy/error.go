package ci_lint_api_proxy

import (
	"encoding/json"
	"net/http"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/logging"
)

type errResponse struct {
	Message string `json:"message"`
}

func sendErr(w http.ResponseWriter, status int, msg string) {
	jsonRes, err := json.Marshal(errResponse{Message: msg})
	if err != nil {
		logging.Error("Failed to serialize error message", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(jsonRes)
}
