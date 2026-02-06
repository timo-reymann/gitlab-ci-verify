package verifier

import (
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/cli"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/git"
	ci_yaml "github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/ci-yaml"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/checks"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/formatter"
)

type GitlabCIVerifier struct {
	configuration *cli.Configuration
	projectRoot   string
	formatter     formatter.FindingsFormatter
}

// NewGitlabCIVerifier creates a new verifier instance
func NewGitlabCIVerifier(c *cli.Configuration, projectRoot string) *GitlabCIVerifier {
	return &GitlabCIVerifier{
		configuration: c,
		projectRoot:   projectRoot,
	}
}

// SetupRego registers all policies and bundles
func (gcv *GitlabCIVerifier) SetupRego() {
	if checks.HasProjectPoliciesOnDisk(gcv.projectRoot) {
		logging.Debug("register project policies")
		checks.RegisterProjectPolicies(gcv.projectRoot)
	}

	if len(gcv.configuration.IncludedOPABundles) > 0 {
		logging.Debug("register opa bundles")
		checks.RegisterRemoteOPABundleChecks(gcv.configuration.IncludedOPABundles)
	}
}

// SetupFormatter sets up the formatter and writes the header
func (gcv *GitlabCIVerifier) SetupFormatter(writer io.Writer, formatterName string) error {
	if formatterName != "" {
		gcv.configuration.OutputFormat = formatterName
	}

	findingsFormatter, err := formatter.Get(gcv.configuration.OutputFormat)
	if err != nil {
		return err
	}
	gcv.formatter = findingsFormatter

	severity := checks.SeverityNameToLevel(gcv.configuration.FailSeverity)
	if severity == -1 {
		return fmt.Errorf("invalid severity level %s", gcv.configuration.FailSeverity)
	}

	err = findingsFormatter.Init(writer)
	if err != nil {
		return err
	}

	return findingsFormatter.Start()
}

// RunChecks runs all checks and prints the results
func (gcv *GitlabCIVerifier) RunChecks(checkInput *checks.CheckInput, checksToRun []checks.Check, failSeverity int, errorHandler func(err error)) bool {
	checkResultChans := checks.RunChecksInParallel(checksToRun, checkInput, func(err error) {
		errorHandler(err)
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
	showedFindings := make([]string, 0)
	for _, finding := range findings {
		if finding.HasCodeIn(gcv.configuration.ExcludedChecks) {
			continue
		}

		ignoredViaComments := checkInput.VirtualCiYaml.GetIgnoredCodes(finding.Line)
		if finding.HasCodeIn(ignoredViaComments) {
			continue
		}

		if slices.Contains(showedFindings, finding.Fingerprint()) {
			continue
		}

		if !shouldFail && finding.HasEqualOrHigherSeverityThan(failSeverity) {
			shouldFail = true
		}

		err := gcv.formatter.Print(&finding)
		errorHandler(err)
		showedFindings = append(showedFindings, finding.Fingerprint())
	}

	err := gcv.formatter.End()
	errorHandler(err)

	return shouldFail
}

// CreateCheckInput creates the check input for running the checks
func (gcv *GitlabCIVerifier) CreateCheckInput() (*checks.CheckInput, error) {
	var err error

	logging.Verbose("read gitlab ci file ", gcv.configuration.GitLabCiFile)
	var ciYamlContent []byte
	if gcv.configuration.GitLabCiFile == "-" {
		ciYamlContent, err = cli.ReadStdinPipe()
	} else {
		ciYamlContent, err = os.ReadFile(gcv.configuration.GitLabCiFile)
	}
	if err != nil {
		return nil, err
	}

	logging.Verbose("load and parse YAML")
	ciYaml, err := ci_yaml.NewCiYamlFile(ciYamlContent)
	if err != nil {
		return nil, err
	}

	var lintRes *ci_yaml.VerificationResultWithRemoteInfo
	var mergedCiYaml *ci_yaml.CiYamlFile

	virtual, err := ci_yaml.CreateVirtualCiYamlFile(gcv.projectRoot, gcv.configuration.GitLabCiFile, ciYaml)
	if err != nil {
		return nil, err
	}
	logging.Debug("Created virtual ci YAML", fmt.Sprintf("\n%s", virtual.Combined.FileContent))

	if gcv.shouldCheckAgainstLintAPI() {
		logging.Debug("Checking against lint API")
		lintRes, err := gcv.checkAgainstLintAPI(lintRes, virtual, mergedCiYaml)
		if err != nil {
			return nil, err
		}

		if lintRes.LintResult.Valid {
			mergedCiYaml, err = ci_yaml.NewCiYamlFile([]byte(lintRes.LintResult.MergedYaml))
			if err != nil {
				return nil, err
			}
		}

	} else {
		logging.Debug("Skipping lint API check")
	}

	return &checks.CheckInput{
		VirtualCiYaml: virtual,
		Configuration: gcv.configuration,
		LintAPIResult: lintRes,
		MergedCiYaml:  mergedCiYaml,
	}, nil
}

func (gcv *GitlabCIVerifier) shouldCheckAgainstLintAPI() bool {
	return !gcv.configuration.Offline && (!gcv.configuration.IsCIEnv() || gcv.configuration.IsCIEnv() && !gcv.configuration.NoLintAPICallInCi)
}

func (gcv *GitlabCIVerifier) checkAgainstLintAPI(lintRes *ci_yaml.VerificationResultWithRemoteInfo, virtual *ci_yaml.VirtualCiYamlFile, mergedCiYaml *ci_yaml.CiYamlFile) (*ci_yaml.VerificationResultWithRemoteInfo, error) {
	logging.Verbose("get remote urls")
	remoteUrls, err := git.GetRemoteUrls(gcv.projectRoot)
	if err != nil {
		return nil, err
	}

	logging.Verbose("parse remote url contents")
	remoteInfos := git.FilterGitlabRemoteUrls(remoteUrls)

	logging.Verbose("validate ci file against gitlab api")
	lintRes, err = ci_yaml.GetFirstValidationResult(&ci_yaml.ValidationResultInput{
		RemoteInfos:      remoteInfos,
		Token:            gcv.configuration.GitlabToken,
		BaseUrlOverwrite: gcv.configuration.GitlabBaseUrlOverwrite(),
		CiYaml:           virtual.Combined.FileContent,
		Timeout:          3 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return lintRes, nil
}
