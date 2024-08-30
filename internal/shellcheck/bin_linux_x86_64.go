//go:build linux && amd64

package shellcheck

import _ "embed"

//go:embed bin/linux-x86_64
var shellcheckBinary []byte
