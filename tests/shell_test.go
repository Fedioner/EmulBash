package tests

import (
	"strings"
	"testing"

	"github.com/Fedioner/EmulBash/src/shell"
)

func TestStubs(t *testing.T) {
	sh := shell.New()
	out, err := sh.Execute("ls -l /tmp")
	if err != nil {
		t.Fatal(err)
	}
	if out != `ls: args ["-l" "/tmp"]` {
		t.Errorf("ls output = %q", out)
	}
	out, err = sh.Execute("cd 'some dir'")
	if err != nil {
		t.Fatal(err)
	}
	if out != `cd: args ["some dir"]` {
		t.Errorf("cd output = %q", out)
	}
}

func TestUnknownCommand(t *testing.T) {
	sh := shell.New()
	_, err := sh.Execute("abc 1 2")
	if err == nil || err.Error() != "abc: command not found" {
		t.Errorf("want command not found, got %v", err)
	}
}

func TestEmptyLine(t *testing.T) {
	sh := shell.New()
	out, err := sh.Execute("   ")
	if out != "" || err != nil {
		t.Errorf("empty line: %q %v", out, err)
	}
}

func TestParseErrorInExecute(t *testing.T) {
	sh := shell.New()
	if _, err := sh.Execute(`ls "abc`); err == nil {
		t.Error("want error for unclosed quote")
	}
}

func TestExit(t *testing.T) {
	sh := shell.New()
	if _, err := sh.Execute("exit"); err != nil {
		t.Fatal(err)
	}
	if !sh.Exited || sh.ExitCode != 0 {
		t.Errorf("exit: exited=%v code=%d", sh.Exited, sh.ExitCode)
	}
}

func TestExitCode(t *testing.T) {
	sh := shell.New()
	if _, err := sh.Execute("exit 3"); err != nil {
		t.Fatal(err)
	}
	if !sh.Exited || sh.ExitCode != 3 {
		t.Errorf("exit 3: exited=%v code=%d", sh.Exited, sh.ExitCode)
	}
}

func TestExitErrors(t *testing.T) {
	sh := shell.New()
	if _, err := sh.Execute("exit abc"); err == nil {
		t.Error("exit abc: want error")
	}
	if _, err := sh.Execute("exit 1 2"); err == nil {
		t.Error("exit 1 2: want error")
	}
	if sh.Exited {
		t.Error("shell should not exit after error")
	}
}

func TestRunShowsPrompt(t *testing.T) {
	sh := shell.New()
	var lines []string
	sh.Run("ls a", func(s string) { lines = append(lines, s) })
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %q", lines)
	}
	if lines[0] != sh.Prompt()+"ls a" {
		t.Errorf("first line = %q", lines[0])
	}
}

func TestTitle(t *testing.T) {
	sh := shell.New()
	want := "Эмулятор - [" + sh.User + "@" + sh.Host + "]"
	if sh.Title() != want {
		t.Errorf("title = %q, want %q", sh.Title(), want)
	}
	if sh.User == "" || strings.Contains(sh.Host, ".") {
		t.Errorf("bad user or host: %q %q", sh.User, sh.Host)
	}
}
