package cmd

import (
	"errors"
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

	fmt, err := formatter.Get(c.OutputFormat)
	handleErr(err, c)

	logging.Verbose("read gitlab ci file ", c.GitLabCiFile)
	ciYamlContent, err := os.ReadFile(c.GitLabCiFile)
	handleErr(err, c)

	ciYaml, err := checks.NewCiYaml(ciYamlContent)
	handleErr(err, c)

	checkInput := checks.CheckInput{
		CiYaml:        ciYaml,
		Configuration: c,
	}

	err = fmt.Init(os.Stdout)
	handleErr(err, c)

	err = fmt.Start()
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

	for _, finding := range findings {
		err := fmt.Print(&finding)
		handleErr(err, c)
	}

	err = fmt.End()
	handleErr(err, c)
}
