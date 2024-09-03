package cmd

import (
	"errors"
	"fmt"
	format_conversion "github.com/timo-reymann/gitlab-ci-verify/internal/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	_ "github.com/timo-reymann/gitlab-ci-verify/internal/shellcheck"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	_ "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
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

	ciYamlDoc, err := format_conversion.ParseYamlNode(ciYamlContent)
	errCheck(err, c)

	var ciYamlParsed map[string]any
	err = ciYamlDoc.Decode(&ciYamlParsed)
	errCheck(err, c)

	checkInput := checks.CheckInput{
		CiYaml: &checks.CiYaml{
			FileContent:   ciYamlContent,
			ParsedYamlMap: ciYamlParsed,
			ParsedYamlDoc: ciYamlDoc,
		},
		Configuration: c,
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	_, _ = fmt.Fprintln(w, strings.Join([]string{
		"Severity",
		"Code",
		"Line",
		"Description",
		"Link",
	}, "\t"))

	for _, check := range checks.AllChecks() {
		results, err := check.Run(&checkInput)
		if err != nil {
			errCheck(err, c)
		}

		for _, result := range results {
			_, _ = fmt.Fprintln(w, strings.Join([]string{
				strings.ToUpper(result.SeverityName()),
				result.Code,
				strconv.Itoa(result.Line),
				result.Message,
				result.Link,
			}, "\t"))
		}
	}
	println("---")
	_ = w.Flush()
}
