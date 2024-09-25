package shellcheck

var ignoredChecksFlags []string

var ignoredChecks = []string{
	"2034",
	"1091",
}

func init() {
	ignoredChecksFlags = make([]string, len(ignoredChecks)*2)
	idx := 0
	for _, ignoredCheck := range ignoredChecks {
		ignoredChecksFlags[idx] = "-e"
		ignoredChecksFlags[idx+1] = "SC" + ignoredCheck
		idx += 2
	}
}
