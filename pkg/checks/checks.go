package checks

var checks = make([]Check, 0)

// register a check
func register(c Check) {
	checks = append(checks, c)
}

// AllChecks available
func AllChecks() []Check {
	return checks
}

func init() {
	register(ShellScriptCheck{})
	register(PipelineLintApiCheck{})
	register(NewGitlabPagesJobCheck())
}
