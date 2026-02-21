package main

import (
	_ "embed"

	"github.com/timo-reymann/gitlab-ci-verify/v2/cmd"
)

//go:embed NOTICE
var noticeContent string

func main() {
	cmd.Execute(noticeContent)
}
