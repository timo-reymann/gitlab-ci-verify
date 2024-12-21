package func_call_bundle

import rego.v1

findings contains finding if {
	true

    finding := {"description": say_hello("World")}
}