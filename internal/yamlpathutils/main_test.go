package yamlpathutils

import (
	"errors"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
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

	path := MustPath(yamlpath.NewPath("."))
	if path == nil {
		t.Fatal("path is not returned when no error is given")
	}
}
