package yamlpathutils

import (
	"errors"
	"testing"
)

func Test_MustPath(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("Expected error to panic")
		}
	}()
	MustPath(nil, errors.New("alarm"))
}
