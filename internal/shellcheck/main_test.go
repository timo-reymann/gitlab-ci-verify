package shellcheck

import "testing"

func TestValidScript(t *testing.T) {
	sc, err := NewShellChecker()
	if err != nil {
		t.Fatal(err)
	}
	res, err := sc.AnalyzeFile("testdata/valid_script.sh")
	if err != nil {
		t.Fatal(err)
	}

	if res.ExitCode != 0 {
		t.Fatalf("Exptected zero exit code but got %d", res.ExitCode)
	}
}

func TestInValidScript(t *testing.T) {
	sc, err := NewShellChecker()
	if err != nil {
		t.Fatal(err)
	}
	res, err := sc.AnalyzeFile("testdata/invalid_script.sh")
	if err != nil {
		t.Fatal(err)
	}

	if res.ExitCode != 1 {
		t.Fatalf("Exptected failure exit code but got %d", res.ExitCode)
	}
}
