package cmd

import (
	"errors"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/cli"
	git2 "github.com/timo-reymann/gitlab-ci-verify/internal/git"
	"github.com/timo-reymann/gitlab-ci-verify/internal/gitlab/ci-yaml"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/formatter"
	"os"
	"slices"
	"time"
)

func handleErr(err error) {
	if errors.Is(err, cli.ErrAbort) {
		os.Exit(0)
	}

	if err != nil {
		println("ERR: " + err.Error())
		os.Exit(2)
	}
}

func Execute() {
	logging.Verbose("create and parse configuration")
	c := cli.NewConfiguration()
	handleErr(c.Parse())

	logging.Verbose("get current working directory")
	projectRoot, err := os.Getwd()
	handleErr(err)

	if checks.HasProjectPoliciesOnDisk(projectRoot) {
		logging.Debug("register project policies")
		checks.RegisterProjectPolicies(projectRoot)
	}

	if len(c.IncludedOPABundles) > 0 {
		logging.Debug("register opa bundles")
		checks.RegisterRemoteOPABundleChecks(c.IncludedOPABundles)
	}

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

	err = findingsFormatter.Start()
	handleErr(err)

	checkInput, err := setupCheckInput(c, projectRoot)
	handleErr(err)
	runChecks(checkInput, c, severity, findingsFormatter)
}

func runChecks(checkInput *checks.CheckInput, c *cli.Configuration, failSeverity int, findingsFormatter formatter.FindingsFormatter) {
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
		if finding.HasCodeIn(c.ExcludedChecks) {
			continue
		}

		ignoredViaComments := checkInput.VirtualCiYaml.GetIgnoredCodes(finding.Line)
		if finding.HasCodeIn(ignoredViaComments) {
			continue
		}

		if !shouldFail && finding.HasEqualOrHigherSeverityThan(failSeverity) {
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

func setupCheckInput(c *cli.Configuration, pwd string) (*checks.CheckInput, error) {
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
	ciYaml, err := ci_yaml.NewCiYamlFile(ciYamlContent)
	handleErr(err)

	var lintRes *ci_yaml.VerificationResultWithRemoteInfo
	var mergedCiYaml *ci_yaml.CiYamlFile

	virtual, err := ci_yaml.CreateVirtualCiYamlFile(pwd, c.GitLabCiFile, ciYaml)
	if err != nil {
		return nil, err
	}
	logging.Debug("Created virtual ci YAML", fmt.Sprintf("\n%s", virtual.Combined.FileContent))

	if !c.IsCIEnv() || c.IsCIEnv() && !c.NoLintAPICallInCi {
		logging.Verbose("get remote urls")
		remoteUrls, err := git2.GetRemoteUrls(pwd)
		handleErr(err)

		logging.Verbose("parse remote url contents")
		remoteInfos := git2.FilterGitlabRemoteUrls(remoteUrls)

		logging.Verbose("validate ci file against gitlab api")
		lintRes, err = ci_yaml.GetFirstValidationResult(&ci_yaml.ValidationResultInput{
			RemoteInfos:      remoteInfos,
			Token:            c.GitlabToken,
			BaseUrlOverwrite: c.GitlabBaseUrlOverwrite(),
			CiYaml:           virtual.Combined.FileContent,
			Timeout:          3 * time.Second,
		})
		handleErr(err)

		mergedCiYaml, err = ci_yaml.NewCiYamlFile([]byte(lintRes.LintResult.MergedYaml))
		handleErr(err)
	}

	checkInput := &checks.CheckInput{
		VirtualCiYaml: virtual,
		Configuration: c,
		LintAPIResult: lintRes,
		MergedCiYaml:  mergedCiYaml,
	}
	return checkInput, err
}
