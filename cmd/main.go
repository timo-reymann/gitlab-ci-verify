package cmd

import (
	"errors"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	_ "github.com/timo-reymann/gitlab-ci-verify/internal/shellcheck"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/formatter"
	_ "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"log"
	"os"
	"slices"
)

func handleErr(err error, c *cli.Configuration) {
	if errors.Is(err, cli.ErrAbort) {
		os.Exit(0)
	}

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
}

func Execute() {
	logging.Verbose("create and parse configuration")
	c := cli.NewConfiguration()
	handleErr(c.Parse(), c)

	envOutputFormat := os.Getenv("GITLAB_CI_VERIFY_OUTPUT_FORMAT")
	if envOutputFormat != "" {
		c.OutputFormat = envOutputFormat
	}

	findingsFormatter, err := formatter.Get(c.OutputFormat)
	handleErr(err, c)

	severity := checks.SeverityNameToLevel(c.FailSeverity)
	if severity == -1 {
		handleErr(fmt.Errorf("invalid severity level %s", c.FailSeverity), c)
	}

	logging.Verbose("read gitlab ci file ", c.GitLabCiFile)
	ciYamlContent, err := os.ReadFile(c.GitLabCiFile)
	handleErr(err, c)

	ciYaml, err := checks.NewCiYaml(ciYamlContent)
	handleErr(err, c)

	checkInput := checks.CheckInput{
		CiYaml:        ciYaml,
		Configuration: c,
	}

	err = findingsFormatter.Init(os.Stdout)
	handleErr(err, c)

	err = findingsFormatter.Start()
	handleErr(err, c)

	checkResultChans := checks.RunChecksInParallel(checks.AllChecks(), checkInput, func(err error) {
		handleErr(err, c)
	})
	findings := make([]checks.CheckFinding, 0)

	for checkResultFindingChan := range checkResultChans {
		for _, finding := range checkResultFindingChan {
			findings = append(findings, finding)
		}
	}

	slices.SortStableFunc(findings, func(a, b checks.CheckFinding) int {
		return a.Compare(b)
	})

	shouldFail := false
	for _, finding := range findings {
		if slices.Contains(c.ExcludedChecks, finding.Code) {
			continue
		}

		if severity >= finding.Severity {
			shouldFail = true
		}

		err := findingsFormatter.Print(&finding)
		handleErr(err, c)
	}

	err = findingsFormatter.End()
	handleErr(err, c)

	if shouldFail {
		os.Exit(1)
	}
}
