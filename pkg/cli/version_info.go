package cli

import (
	"errors"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/shellcheck"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/buildinfo"
	"os"
	"runtime"
	"text/tabwriter"
)

// ErrAbort states that the regular execution should be aborted
var ErrAbort = errors.New("abort")

func addLine(w *tabwriter.Writer, heading string, val string) {
	_, _ = fmt.Fprintf(w, heading+"\t%s\n", val)
}

// PrintVersionInfo prints a tabular list with build info
func PrintVersionInfo() {
	PrintCompactInfo()
	println()
	println("Build information")
	w := tabwriter.NewWriter(os.Stderr, 10, 1, 10, byte(' '), tabwriter.TabIndent)
	addLine(w, "GitSha", buildinfo.GitSha)
	addLine(w, "Version", buildinfo.Version)
	addLine(w, "BuildTime", buildinfo.BuildTime)
	addLine(w, "ShellCheck-Version", shellCheckInfo())
	addLine(w, "Go-Version", runtime.Version())
	addLine(w, "OS/Arch", runtime.GOOS+"/"+runtime.GOARCH)
	_ = w.Flush()
}

func shellCheckInfo() string {
	sc, _ := shellcheck.NewShellChecker()
	return sc.Version()
}

// PrintCompactInfo outputs a online header with basic infos
func PrintCompactInfo() {
	fmt.Printf("gitlab-ci-verify %s (%s) by Timo Reymann\n", buildinfo.Version, buildinfo.BuildTime)
}
