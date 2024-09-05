package formatter

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"testing"
)

func verifyFormatter(t *testing.T, f FindingsFormatter, findings []*checks.CheckFinding, expectedOutput []byte) {
	buf := bytes.NewBuffer([]byte{})

	if err := f.Init(buf); err != nil {
		t.Fatal(err)
	}

	if err := f.Start(); err != nil {
		t.Fatal(err)
	}

	for _, finding := range findings {
		err := f.Print(finding)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := f.End(); err != nil {
		t.Fatal(err)
	}

	out := buf.Bytes()
	if !cmp.Equal(expectedOutput, out) {
		t.Fatalf("Expected output:\n%s\ngot:\n%s", string(expectedOutput), string(out))
	}
}
