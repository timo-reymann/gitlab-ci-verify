package cli

import (
	flag "github.com/spf13/pflag"
)

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
	conf.addStringFlag(&conf.GitlabBaseUrl, "gitlab-base-url", "", "auto-detect", "Set the gitlab base url explicitly in case detection does not work or your clone and base url differs")
	conf.addStringFlag(&conf.GitlabToken, "gitlab-token", "", "", "Gitlab token to use")
}

func (conf *Configuration) Help() {
	PrintCompactInfo()
	println("gitlab-ci-verify [-options]")
	flag.PrintDefaults()
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
		PrintVersionInfo()
		return ErrAbort
	}

	return nil
}

// NewConfiguration creates a new configuration
func NewConfiguration() *Configuration {
	return &Configuration{}
}