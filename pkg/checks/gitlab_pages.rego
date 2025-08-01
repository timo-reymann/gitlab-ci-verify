package gitlab_ci_verify_gitlab_pages_check

import data.gitlab_ci_verify
import rego.v1

findings contains finding if {
    artifact_paths := input.mergedYaml.pages.artifacts.paths

	count([artifact_paths |
		some artifact_path in artifact_paths
		startswith(artifact_path, "public")
	]) == 0

	finding := gitlab_ci_verify.warning_with_link(
	    "GL-201",
	    "Pages job does not contain artifacts with a public path.",
	     yamlPathToLineNumber(".pages.artifacts.paths"),
	     "https://gitlab-ci-verify.timo-reymann.de/findings/GL-202.html"
	)
}

findings contains finding if {
    input.mergedYaml.pages.artifacts == {}

	finding := gitlab_ci_verify.warning_with_link(
	    "GL-202",
	    "Pages job does not define artifacts.",
	     yamlPathToLineNumber(".pages.artifacts"),
	     "https://gitlab-ci-verify.timo-reymann.de/findings/GL-202.html"
	)
}
