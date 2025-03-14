package location

import (
	"fmt"
	"path/filepath"
)

type Location struct {
	File string
	Line int
}

func NewLocation(file string, line int) *Location {
	return &Location{
		File: file,
		Line: line,
	}
}

func (l *Location) Absolute() (*Location, error) {
	absPath, err := filepath.Abs(l.File)
	if err != nil {
		return nil, err
	}

	return NewLocation(absPath, l.Line), nil
}

func (l *Location) String() string {
	return fmt.Sprintf("%s:%d", l.File, l.Line)
}
