//go:build darwin && arm64

package shellcheck

import _ "embed"

//go:embed bin/darwin-aarch64
var shellcheckBinary []byte
