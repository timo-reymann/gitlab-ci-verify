package logging

import (
	"io"
	"log"
	"os"
)

const LevelSilent = 0
const LevelDebug = 1
const LevelVerbose = 2

var Level = LevelSilent

type consoleLogger struct {
	verbose *log.Logger
	debug   *log.Logger
	warn    *log.Logger
}

func newConsoleLogger(w io.Writer) *consoleLogger {
	return &consoleLogger{
		debug:   log.New(w, "DEBUG: ", log.Lmsgprefix),
		verbose: log.New(w, "VERBOSE: ", log.Lmsgprefix),
		warn:    log.New(w, "WARN: ", log.Lmsgprefix),
	}
}

var Logger = newConsoleLogger(os.Stderr)

func (l *consoleLogger) Debug(v ...any) {
	if Level < LevelDebug {
		return
	}
	l.debug.Println(v...)
}

func (l *consoleLogger) Verbose(v ...any) {
	if Level < LevelVerbose {
		return
	}
	l.verbose.Println(v...)
}

func (l *consoleLogger) Warn(v ...any) {
	if Level < LevelDebug {
		return
	}
	l.warn.Println(v...)
}
