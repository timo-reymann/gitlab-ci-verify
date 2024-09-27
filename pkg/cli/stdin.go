package cli

import (
	"errors"
	"io"
	"os"
)

func isStdinPiped() bool {
	stat, err := os.Stdin.Stat()
	return err == nil && (stat.Mode()&os.ModeCharDevice) == 0
}

func ReadStdinPipe() ([]byte, error) {
	if isStdinPiped() {
		return io.ReadAll(os.Stdin)
	}

	return nil, errors.New("stdin is not readable")
}
