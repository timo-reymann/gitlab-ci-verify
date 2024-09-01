package logging

import (
	"bytes"
	"fmt"
	"testing"
)

func TestConsoleLogger(t *testing.T) {
	testCases := []struct {
		logLevel       int
		typ            string
		message        string
		expectedOutput string
	}{
		{
			logLevel:       LevelSilent,
			typ:            "debug",
			message:        "This is a debug message",
			expectedOutput: "",
		},
		{
			logLevel:       LevelDebug,
			typ:            "debug",
			message:        "This is a debug message",
			expectedOutput: "DEBUG: This is a debug message\n",
		},
		{
			logLevel:       LevelVerbose,
			typ:            "debug",
			message:        "This is a debug message",
			expectedOutput: "DEBUG: This is a debug message\n",
		},
		{
			logLevel:       LevelSilent,
			typ:            "warn",
			message:        "This is a warning message",
			expectedOutput: "",
		},
		{
			logLevel:       LevelDebug,
			typ:            "warn",
			message:        "This is a warning message",
			expectedOutput: "WARN: This is a warning message\n",
		},
		{
			logLevel:       LevelVerbose,
			typ:            "verbose",
			message:        "This is a verbose message",
			expectedOutput: "VERBOSE: This is a verbose message\n",
		},
		{
			logLevel:       LevelSilent,
			typ:            "verbose",
			message:        "This is a verbose message",
			expectedOutput: "",
		},
		{
			logLevel:       LevelDebug,
			typ:            "verbose",
			message:        "This is a verbose message",
			expectedOutput: "",
		},
		{
			logLevel:       LevelVerbose,
			typ:            "verbose",
			message:        "This is a verbose message",
			expectedOutput: "VERBOSE: This is a verbose message\n",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", tc.logLevel, tc.typ), func(t *testing.T) {
			Level = tc.logLevel
			buffer := bytes.NewBuffer([]byte{})
			logger := newConsoleLogger(buffer)

			switch tc.typ {
			case "debug":
				logger.Debug(tc.message)
			case "verbose":
				logger.Verbose(tc.message)
			case "warn":
				logger.Warn(tc.message)
			}

			if tc.expectedOutput != buffer.String() {
				t.Fatalf("Expected output %s, but got %s", tc.expectedOutput, buffer.String())
			}
		})

	}
}
