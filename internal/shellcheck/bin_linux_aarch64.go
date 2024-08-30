//go:build linux && arm64

package shellcheck

import _ "embed"

//go:embed bin/linux-aarch64
var shellcheckBinary []byte
