package cmd

import (
	"errors"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/git"
	ciyaml "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"log"
	"os"
	"time"
)

func errCheck(err error, c *cli.Configuration) {
	if errors.Is(err, cli.ErrAbort) {
		os.Exit(0)
	}

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
}

func Execute() {
	c := cli.NewConfiguration()
	errCheck(c.Parse(), c)

	ciYamlContent, err := os.ReadFile(c.GitLabCiFile)
	errCheck(err, c)

	pwd, err := os.Getwd()
	errCheck(err, c)

	remoteUrls, err := git.GetRemoteUrls(pwd)
	errCheck(err, c)
	remoteInfos := git.FilterGitlabRemoteUrls(remoteUrls)

	res, err := ciyaml.GetFirstValidationResult(remoteInfos, c.GitlabToken, "", ciYamlContent, 3*time.Second)
	errCheck(err, c)

	if res.LintResult.Valid {
		println("Valid configuration")
		os.Exit(0)
	} else {
		println("Invalid configuration")
		fmt.Printf(" -> warnings: %v\n", res.LintResult.Warnings)
		fmt.Printf(" -> errors:   %v\n", res.LintResult.Errors)
		os.Exit(1)
	}
}
