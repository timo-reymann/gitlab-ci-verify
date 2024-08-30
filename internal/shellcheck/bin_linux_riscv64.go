//go:build linux && riscv64

package shellcheck

import _ "embed"

//go:embed bin/linux-riscv64
var shellcheckBinary []byte
