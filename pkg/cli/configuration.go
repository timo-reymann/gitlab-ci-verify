package cli

import (
	"errors"
	flag "github.com/spf13/pflag"
	"github.com/timo-reymann/gitlab-ci-verify/internal/buildinfo"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"os"
)

// AutoDetectValue indicates that the value should be detected automatically from a set of sources
var AutoDetectValue = "auto-detect"

// ErrAbort states that the regular execution should be aborted
var ErrAbort = errors.New("abort")

// Configuration for the CLI
type Configuration struct {
	GitLabCiFile       string
	GitlabBaseUrl      string
	GitlabToken        string
	Verbose            bool
	Debug              bool
	ShellcheckFlags    string
	OutputFormat       string
	FailSeverity       string
	ExcludedChecks     []string
	NoLintAPICallInCi  bool
	IncludedOPABundles []string
}

func (conf *Configuration) addBoolFlag(field *bool, long string, short string, val bool, usage string) {
	flag.BoolVarP(field, long, short, val, usage)
}

func (conf *Configuration) addStringsFlag(field *[]string, long string, short string, val []string, usage string) {
	if short == "" {
		flag.StringSliceVar(field, long, val, usage)
	} else {
		flag.StringSliceVarP(field, long, short, val, usage)
	}
}

func (conf *Configuration) addStringFlag(field *string, long string, short string, val string, usage string) {
	flag.StringVarP(field, long, short, val, usage)
}

func (conf *Configuration) defineFlags() {
	conf.addStringFlag(&conf.GitLabCiFile, "gitlab-ci-file", "", ".gitlab-ci.yml", "The Yaml file used to configure GitLab CI")
	conf.addStringFlag(&conf.GitlabBaseUrl, "gitlab-base-url", "", AutoDetectValue, "Set the gitlab base url explicitly in case detection does not work or your clone and base url differs")
	conf.addStringFlag(&conf.GitlabToken, "gitlab-token", "", "", "Gitlab token to use, if not specified the netrc is evaluated and if that also does not contain credentials, tries to load the environment variable GITLAB_TOKEN")
	conf.addBoolFlag(&conf.Debug, "debug", "", false, "Enable debug output")
	conf.addBoolFlag(&conf.Verbose, "verbose", "", false, "Enable verbose output")
	conf.addStringFlag(&conf.ShellcheckFlags, "shellcheck-flags", "", "", "Pass custom flags to shellcheck")
	conf.addStringFlag(&conf.OutputFormat, "format", "f", "text", "Format for the output, valid options are json, table and text. If GITLAB_CI_VERIFY_OUTPUT_FORMAT this parameter is ignored")
	conf.addStringFlag(&conf.FailSeverity, "severity", "S", "style", "Set the severity level on which to consider findings as errors and exiting with non zero exit code.")
	conf.addStringsFlag(&conf.ExcludedChecks, "exclude", "E", []string{}, "Exclude the given check codes")
	conf.addBoolFlag(&conf.NoLintAPICallInCi, "no-lint-api-in-ci", "", false, "Add this flag to avoid validating against Pipeline Check API, as its assumed that running in CI is proof enough the syntax is valid. Please note that checks relying on the merged YAML will also not be executed in that case.")
	conf.addStringsFlag(&conf.IncludedOPABundles, "include-opa-bundle", "I", []string{}, "Include remote OPA bundles for checks")
}

func (conf *Configuration) Help() {
	buildinfo.PrintCompactInfo(os.Stdout)
	println("gitlab-ci-verify [-options]")
	flag.PrintDefaults()
}

// GitlabBaseUrlOverwrite returns an empty string in case the value should be automatically detected
// or otherwise the value of the cli parameter
func (conf *Configuration) GitlabBaseUrlOverwrite() string {
	if conf.GitlabBaseUrl == AutoDetectValue {
		return ""
	}

	return conf.GitlabBaseUrl
}

func (conf *Configuration) configureLogLevel() {
	if conf.Verbose {
		logging.Level = logging.LevelVerbose
	} else if conf.Debug {
		logging.Level = logging.LevelDebug
	}
}

// Parse the configuration from cli args
func (conf *Configuration) Parse() error {
	conf.defineFlags()

	isHelp := flag.BoolP("help", "h", false, "Show available commands")
	isVersion := flag.Bool("version", false, "Show version info")
	flag.Parse()
	conf.configureLogLevel()

	if *isHelp {
		conf.Help()
		return ErrAbort
	} else if *isVersion {
		buildinfo.PrintVersionInfo(os.Stderr)
		return ErrAbort
	}

	return nil
}

func (conf *Configuration) IsCIEnv() bool {
	return os.Getenv("CI") != ""
}

// NewConfiguration creates a new configuration
func NewConfiguration() *Configuration {
	return &Configuration{}
}
