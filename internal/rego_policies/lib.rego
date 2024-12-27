package gitlab_ci_verify

import rego.v1

create_finding(code, severity, message, line, link) := {
	"code": code,
	"severity": severity,
	"message": message,
	"link": link,
	"line": line,
}

warning(code, msg, line) := create_finding(code, "WARNING", msg, line, "")

warning_with_link(code, msg, line, link) := create_finding(code, "WARNING", msg, line, link)

error(code, msg, line) := create_finding(code, "ERROR", msg, line, "")

error_with_link(code, msg, line, link) := create_finding(code, "ERROR", msg, line, link)

info(code, msg, line) := create_finding(code, "INFO", msg, line, "")

info_with_link(code, msg, line, link) := create_finding(code, "INFO", msg, line, link)