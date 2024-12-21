package gitlab_ci_verify

import rego.v1

create_finding(code, severity, message, link) := {
	"code": code,
	"severity": severity,
	"message": message,
	"link": link,
	"line": -1,
}

warning(code, msg) := create_finding(code, "WARNING", msg, "")

error(code, msg) := create_finding(code, "ERROR", msg, "")

info(code, msg) := create_finding(code, "INFO", msg, "")
