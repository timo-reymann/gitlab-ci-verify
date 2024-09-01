package cli

import (
	"errors"
	flag "github.com/spf13/pflag"
	"github.com/timo-reymann/gitlab-ci-verify/internal/buildinfo"
	"os"
)

// AutoDetectValue indicates that the value should be detected automatically from a set of sources
var AutoDetectValue = "auto-detect"

// ErrAbort states that the regular execution should be aborted
var ErrAbort = errors.New("abort")

// Configuration for the CLI
type Configuration struct {
	GitLabCiFile  string
	GitlabBaseUrl string
	GitlabToken   string
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
	conf.addStringFlag(&conf.GitlabToken, "gitlab-token", "", "", "Gitlab token to use")
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

// Parse the configuration from cli args
func (conf *Configuration) Parse() error {
	conf.defineFlags()

	isHelp := flag.BoolP("help", "h", false, "Show available commands")
	isVersion := flag.Bool("version", false, "Show version info")
	flag.Parse()

	if *isHelp {
		conf.Help()
		return ErrAbort
	} else if *isVersion {
		buildinfo.PrintVersionInfo(os.Stderr)
		return ErrAbort
	}

	return nil
}

// NewConfiguration creates a new configuration
func NewConfiguration() *Configuration {
	return &Configuration{}
}
