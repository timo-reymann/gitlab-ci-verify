//go:generate sh -c "cp ../../NOTICE NOTICE"

package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/buildinfo"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/api"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/logging"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/ci_lint_api_proxy"
)

func setLogLevelFromEnv() {
	logLevel := os.Getenv("LOG_LEVEL")

	switch logLevel {
	case "verbose":
		logging.Level = logging.LevelVerbose
		break
	case "debug":
		logging.Level = logging.LevelDebug
		break
	default:
		logging.Level = logging.LevelSilent
		break
	}

	if logLevel == "" {
		logging.Level = logging.LevelSilent
	}
}

const endpoint = ":8080"

func main() {
	isVersion := flag.Bool("version", false, "Show version info")
	isLicense := flag.Bool("license", false, "Show license info")
	gitlabBaseUrl := flag.String("gitlab-base-url", "", "Base URL of GitLab instance")
	flag.Parse()

	if *isVersion {
		buildinfo.PrintVersionInfo("gitlab-ci-lint-api-proxy", os.Stdout)
		return
	} else if *isLicense {
		fmt.Println(noticeContent)
		return
	}

	if *gitlabBaseUrl == "" {
		logging.Error("GitLab base URL must be provided via --gitlab-base-url flag")
		os.Exit(1)
	}

	setLogLevelFromEnv()

	gitlabClient := api.NewClientWithMultiTokenSources(*gitlabBaseUrl, os.Getenv("GITLAB_TOKEN"))
	ctx := ci_lint_api_proxy.ApiContext{GitLabClient: gitlabClient}

	logging.Debug("Starting server on ", endpoint)
	err := ci_lint_api_proxy.Serve(endpoint, ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logging.Error("Failed to start HTTP server", err)
		os.Exit(1)
	}
}
