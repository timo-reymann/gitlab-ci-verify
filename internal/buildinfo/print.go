package buildinfo

import (
	"fmt"
	"io"
	"runtime"
	"text/tabwriter"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/shellcheck"
)

func shellCheckInfo() string {
	sc, _ := shellcheck.NewShellChecker()
	defer func() {
		if sc != nil {
			_ = sc.Close()
		}
	}()
	return sc.Version()
}

func addLine(w *tabwriter.Writer, heading string, val string) {
	_, _ = fmt.Fprintf(w, heading+"\t%s\n", val)
}

// PrintCompactInfo with minimal build information
func PrintCompactInfo(binary string, w io.Writer) {
	_, _ = fmt.Fprintf(w, "%s %s (%s) by Timo Reymann\n", binary, Version, BuildTime)
}

// PrintVersionInfo prints a tabular list with build info
func PrintVersionInfo(binary string, w io.Writer) {
	PrintCompactInfo(binary, w)
	_, _ = fmt.Fprint(w, "\nBuild information\n")
	tw := tabwriter.NewWriter(w, 10, 1, 10, byte(' '), tabwriter.TabIndent)
	addLine(tw, "GitSha", GitSha)
	addLine(tw, "Version", Version)
	addLine(tw, "BuildTime", BuildTime)
	addLine(tw, "ShellCheck-Version", shellCheckInfo())
	addLine(tw, "Go-Version", runtime.Version())
	addLine(tw, "OS/Arch", runtime.GOOS+"/"+runtime.GOARCH)
	_ = tw.Flush()
}
