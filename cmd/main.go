package cmd

import (
	"errors"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	_ "github.com/timo-reymann/gitlab-ci-verify/internal/shellcheck"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/formatter"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/git"
	_ "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	ciyaml "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"log"
	"os"
	"slices"
	"time"
)

func handleErr(err error) {
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
	handleErr(c.Parse())

	envOutputFormat := os.Getenv("GITLAB_CI_VERIFY_OUTPUT_FORMAT")
	if envOutputFormat != "" {
		c.OutputFormat = envOutputFormat
	}

	findingsFormatter, err := formatter.Get(c.OutputFormat)
	handleErr(err)

	severity := checks.SeverityNameToLevel(c.FailSeverity)
	if severity == -1 {
		handleErr(fmt.Errorf("invalid severity level %s", c.FailSeverity))
	}

	err = findingsFormatter.Init(os.Stdout)
	handleErr(err)

	err, checkInput := setupCheckInput(c)
	runChecks(checkInput, c, severity, findingsFormatter)

	err = findingsFormatter.Start()
	handleErr(err)
}

func runChecks(checkInput checks.CheckInput, c *cli.Configuration, severity int, findingsFormatter formatter.FindingsFormatter) {
	checkResultChans := checks.RunChecksInParallel(checks.AllChecks(), checkInput, func(err error) {
		handleErr(err)
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
		handleErr(err)
	}

	err := findingsFormatter.End()
	handleErr(err)

	if shouldFail {
		os.Exit(1)
	}
}

func setupCheckInput(c *cli.Configuration) (error, checks.CheckInput) {
	var err error

	logging.Verbose("read gitlab ci file ", c.GitLabCiFile)
	var ciYamlContent []byte
	if c.GitLabCiFile == "-" {
		ciYamlContent, err = cli.ReadStdinPipe()
	} else {
		ciYamlContent, err = os.ReadFile(c.GitLabCiFile)
	}
	handleErr(err)

	logging.Verbose("load and parse YAML")
	ciYaml, err := checks.NewCiYaml(ciYamlContent)
	handleErr(err)

	var lintRes *ciyaml.VerificationResultWithRemoteInfo
	var mergedCiYaml *checks.CiYaml

	if !c.NoLintAPICallInCi {
		logging.Verbose("get current working directory")
		pwd, err := os.Getwd()
		handleErr(err)

		logging.Verbose("get remote urls")
		remoteUrls, err := git.GetRemoteUrls(pwd)
		handleErr(err)

		logging.Verbose("parse remote url contents")
		remoteInfos := git.FilterGitlabRemoteUrls(remoteUrls)

		logging.Verbose("validate ci file against gitlab api")
		lintRes, err = ciyaml.GetFirstValidationResult(remoteInfos, c.GitlabToken, c.GitlabBaseUrlOverwrite(), ciYamlContent, 3*time.Second)
		handleErr(err)

		mergedCiYaml, err = checks.NewCiYaml([]byte(lintRes.LintResult.MergedYaml))
		handleErr(err)
	}

	checkInput := checks.CheckInput{
		CiYaml:        ciYaml,
		Configuration: c,
		LintAPIResult: lintRes,
		MergedCiYaml:  mergedCiYaml,
	}
	return err, checkInput
}
