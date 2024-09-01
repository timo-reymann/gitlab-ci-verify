package cmd

import (
	"errors"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
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
	logging.Verbose("create and parse configuration")
	c := cli.NewConfiguration()
	errCheck(c.Parse(), c)

	logging.Verbose("read gitlab ci file ", c.GitLabCiFile)
	ciYamlContent, err := os.ReadFile(c.GitLabCiFile)
	errCheck(err, c)

	logging.Verbose("get current working directory")
	pwd, err := os.Getwd()
	errCheck(err, c)

	logging.Verbose("get remote urls")
	remoteUrls, err := git.GetRemoteUrls(pwd)
	errCheck(err, c)
	logging.Verbose("parse remote url contents")
	remoteInfos := git.FilterGitlabRemoteUrls(remoteUrls)

	logging.Verbose("validate ci file against gitlab api")
	res, err := ciyaml.GetFirstValidationResult(remoteInfos, c.GitlabToken, c.GitlabBaseUrlOverwrite(), ciYamlContent, 3*time.Second)
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
