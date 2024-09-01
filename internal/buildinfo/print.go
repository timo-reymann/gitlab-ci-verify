package buildinfo

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/shellcheck"
	"io"
	"runtime"
	"text/tabwriter"
)

func shellCheckInfo() string {
	sc, _ := shellcheck.NewShellChecker()
	return sc.Version()
}

func addLine(w *tabwriter.Writer, heading string, val string) {
	_, _ = fmt.Fprintf(w, heading+"\t%s\n", val)
}

// PrintCompactInfo with minimal build information
func PrintCompactInfo(w io.Writer) {
	_, _ = fmt.Fprintf(w, "gitlab-ci-verify %s (%s) by Timo Reymann\n", Version, BuildTime)
}

// PrintVersionInfo prints a tabular list with build info
func PrintVersionInfo(w io.Writer) {
	PrintCompactInfo(w)
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
