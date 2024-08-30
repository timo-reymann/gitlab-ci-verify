//go:build windows && amd64

package shellcheck

import _ "embed"

//go:embed bin/windows.exe
var shellcheckBinary []byte
