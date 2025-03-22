package cmd

import (
	"errors"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/cli"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/verifier"
	"os"
)

func handleErr(err error) {
	if errors.Is(err, cli.ErrAbort) {
		os.Exit(0)
	}

	if err != nil {
		println("ERR: " + err.Error())
		os.Exit(2)
	}
}

// Execute runs the verifier with the given configuration
func Execute() {
	logging.Verbose("create and parse configuration")
	c := cli.NewConfiguration()
	handleErr(c.Parse())

	logging.Verbose("get current working directory")
	projectRoot, err := os.Getwd()
	handleErr(err)

	gcv := verifier.NewGitlabCIVerifier(c, projectRoot)
	gcv.SetupRego()

	failSeverity := checks.SeverityNameToLevel(c.FailSeverity)
	if failSeverity == -1 {
		handleErr(fmt.Errorf("invalid severity level %s", c.FailSeverity))
	}

	err = gcv.SetupFormatter(os.Stdout, os.Getenv("GITLAB_CI_VERIFY_OUTPUT_FORMAT"))
	if err != nil {
		handleErr(err)
	}

	checkInput, err := gcv.CreateCheckInput()
	handleErr(err)

	gcv.RunChecks(checkInput, failSeverity, handleErr)
}
