package valid_bundle

import data.nonexistent
import rego.v1

findings contains finding if {
	nonexistent.do_things()
}