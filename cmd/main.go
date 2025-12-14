package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/cli"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/checks"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/verifier"
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

func getOutputWriter(outputFile string) (io.Writer, func(), error) {
	if outputFile == "" {
		return os.Stdout, func() {}, nil
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return file, func() {
		file.Close()
	}, nil
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

	// Get output writer (file or stdout)
	writer, cleanup, err := getOutputWriter(c.OutputFile)
	if err != nil {
		handleErr(err)
	}
	defer cleanup()

	err = gcv.SetupFormatter(writer, os.Getenv("GITLAB_CI_VERIFY_OUTPUT_FORMAT"))
	if err != nil {
		handleErr(err)
	}

	checkInput, err := gcv.CreateCheckInput()
	handleErr(err)

	shouldFail := gcv.RunChecks(checkInput, checks.AllChecks(), failSeverity, handleErr)
	if shouldFail {
		os.Exit(1)
	}
}
