package rego_policies

import (
	"github.com/open-policy-agent/opa/topdown/print"
	"testing"
)

func TestLogPrinter_Print(t *testing.T) {
	prnt := logPrinter{}
	prnt.Print(print.Context{}, "foo")
}
