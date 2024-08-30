//go:build darwin && amd64

package shellcheck

import _ "embed"

//go:embed bin/darwin-x86_64
var shellcheckBinary []byte
