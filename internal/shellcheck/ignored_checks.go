package shellcheck

import "fmt"

var ignoredChecksFlags []string

var ignoredChecks = []string{
	"2034",
}

func init() {
	ignoredChecksFlags = make([]string, len(ignoredChecks)*2)
	for idx, ignoredCheck := range ignoredChecks {
		ignoredChecksFlags[idx] = "-e"
		ignoredChecksFlags[idx+1] = fmt.Sprintf("SC%s", ignoredCheck)
	}
}
