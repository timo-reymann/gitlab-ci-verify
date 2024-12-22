package valid_bundle

import data.gitlab_ci_verify
import rego.v1

findings contains finding if {
	true

	finding := gitlab_ci_verify.warning("420", "always triggers")
}