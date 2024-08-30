package cmd

import (
	"errors"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	"log"
	"os"
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
}
