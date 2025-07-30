package rego_policies

import (
	"github.com/open-policy-agent/opa/topdown/print"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/logging"
)

type logPrinter struct {
}

func (p logPrinter) Print(_ print.Context, s string) error {
	logging.Debug(s)
	return nil
}
