package tests

import (
	"strings"
	"testing"

	"github.com/Fedioner/EmulBash/src/shell"
)

func runScript(t *testing.T, sh *shell.Shell, text string) ([]string, error) {
	t.Helper()
	path := writeFile(t, t.TempDir(), "script.txt", text)
	var lines []string
	err := sh.RunScript(path, func(s string) { lines = append(lines, s) })
	return lines, err
}

func TestRunScript(t *testing.T) {
	sh := shell.New()
	lines, err := runScript(t, sh, "# comment\n\nls a\ncd b\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		sh.Prompt() + "ls a", `ls: args ["a"]`,
		sh.Prompt() + "cd b", `cd: args ["b"]`,
	}
	if strings.Join(lines, "\n") != strings.Join(want, "\n") {
		t.Errorf("got %q, want %q", lines, want)
	}
}

func TestRunScriptStopsOnError(t *testing.T) {
	sh := shell.New()
	lines, err := runScript(t, sh, "ls\nbad\nexit\n")
	if err == nil || !strings.Contains(err.Error(), "строке 2") {
		t.Fatalf("want error on line 2, got %v", err)
	}
	if sh.Exited {
		t.Error("exit after error should not run")
	}
	if lines[len(lines)-1] != "bad: command not found" {
		t.Errorf("last line = %q", lines[len(lines)-1])
	}
}

func TestRunScriptExit(t *testing.T) {
	sh := shell.New()
	lines, err := runScript(t, sh, "exit 5\nls\n")
	if err != nil {
		t.Fatal(err)
	}
	if !sh.Exited || sh.ExitCode != 5 {
		t.Errorf("exited=%v code=%d", sh.Exited, sh.ExitCode)
	}
	if len(lines) != 1 {
		t.Errorf("commands after exit should not run: %q", lines)
	}
}

func TestRunScriptMissing(t *testing.T) {
	sh := shell.New()
	if err := sh.RunScript("/no/such/script.txt", func(string) {}); err == nil {
		t.Error("want error for missing script")
	}
}
