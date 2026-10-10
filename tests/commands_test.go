package tests

import (
	"strings"
	"testing"

	"github.com/Fedioner/EmulBash/src/shell"
	"github.com/Fedioner/EmulBash/src/vfs"
)

func newTestShell(t *testing.T) *shell.Shell {
	t.Helper()
	path := makeZip(t, []zipEntry{
		{"home/user/docs/work.txt", "line1\nline2\n"},
		{"home/user/notes.txt", "hello"},
		{"home/user/.hidden", "secret"},
		{"etc/hosts", "127.0.0.1 localhost\n"},
		{"etc/привет.txt", "абв\n"},
		{"tmp/", ""},
	})
	v, err := vfs.LoadZip(path)
	if err != nil {
		t.Fatal(err)
	}
	sh := shell.New()
	sh.FS = v
	return sh
}

func checkOutput(t *testing.T, sh *shell.Shell, line, want string) {
	t.Helper()
	out, err := sh.Execute(line)
	if err != nil {
		t.Errorf("%s: unexpected error %v", line, err)
		return
	}
	if out != want {
		t.Errorf("%s:\ngot  %q\nwant %q", line, out, want)
	}
}

func checkError(t *testing.T, sh *shell.Shell, line, want string) {
	t.Helper()
	_, err := sh.Execute(line)
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("%s: want error with %q, got %v", line, want, err)
	}
}

func TestLs(t *testing.T) {
	sh := newTestShell(t)
	checkOutput(t, sh, "ls", "etc/\nhome/\ntmp/")
	checkOutput(t, sh, "ls /home/user", "docs/\nnotes.txt")
	checkOutput(t, sh, "ls -a /home/user", ".hidden\ndocs/\nnotes.txt")
	checkOutput(t, sh, "ls /etc/hosts", "/etc/hosts")
	checkOutput(t, sh, "ls /tmp", "")
	checkOutput(t, sh, "ls /etc /tmp", "/etc:\nhosts\nпривет.txt\n\n/tmp:")
}

func TestLsLong(t *testing.T) {
	sh := newTestShell(t)
	want := "drwxr-xr-x      1 docs/\n-rw-r--r--      5 notes.txt"
	checkOutput(t, sh, "ls -l /home/user", want)
	want = "-rw-r--r--      6 .hidden\n" + want
	checkOutput(t, sh, "ls -la /home/user", want)
	checkOutput(t, sh, "ls -a -l /home/user", want)
}

func TestLsErrors(t *testing.T) {
	sh := newTestShell(t)
	checkError(t, sh, "ls /nope", "cannot access '/nope': no such file or directory")
	checkError(t, sh, "ls /etc/hosts/x", "not a directory")
	checkError(t, sh, "ls -x", "invalid option -- 'x'")
	out, err := sh.Execute("ls /nope /tmp /etc")
	if err == nil || out != "/tmp:\n\n/etc:\nhosts\nпривет.txt" {
		t.Errorf("ls with bad and good paths: %q %v", out, err)
	}
}

func TestCd(t *testing.T) {
	sh := newTestShell(t)
	steps := []struct{ cmd, cwd string }{
		{"cd /home/user", "/home/user"},
		{"cd docs", "/home/user/docs"},
		{"cd ..", "/home/user"},
		{"cd ./docs/../../user/", "/home/user"},
		{"cd", "/"},
		{"cd ../..", "/"},
		{"cd etc", "/etc"},
	}
	for _, s := range steps {
		checkOutput(t, sh, s.cmd, "")
		if sh.FS.Cwd != s.cwd {
			t.Errorf("%s: cwd = %q, want %q", s.cmd, sh.FS.Cwd, s.cwd)
		}
	}
	if !strings.HasSuffix(sh.Prompt(), ":/etc$ ") {
		t.Errorf("prompt = %q", sh.Prompt())
	}
	checkOutput(t, sh, "ls", "hosts\nпривет.txt")
}

func TestCdErrors(t *testing.T) {
	sh := newTestShell(t)
	checkError(t, sh, "cd /nope", "cd: /nope: no such file or directory")
	checkError(t, sh, "cd /etc/hosts", "cd: /etc/hosts: not a directory")
	checkError(t, sh, "cd /etc /tmp", "cd: too many arguments")
	if sh.FS.Cwd != "/" {
		t.Errorf("cwd changed after error: %q", sh.FS.Cwd)
	}
}

func TestRev(t *testing.T) {
	sh := newTestShell(t)
	checkOutput(t, sh, "rev /home/user/docs/work.txt", "1enil\n2enil")
	checkOutput(t, sh, "rev /etc/привет.txt", "вба")
	checkOutput(t, sh, "cd /home/user", "")
	checkOutput(t, sh, "rev notes.txt docs/work.txt", "olleh\n1enil\n2enil")
	checkOutput(t, sh, "rev /tmp/../home/user/notes.txt", "olleh")
}

func TestRevErrors(t *testing.T) {
	sh := newTestShell(t)
	checkError(t, sh, "rev", "rev: missing file operand")
	checkError(t, sh, "rev /nope.txt", "rev: /nope.txt: no such file or directory")
	checkError(t, sh, "rev /home", "rev: /home: is a directory")
}

func TestFind(t *testing.T) {
	sh := newTestShell(t)
	checkOutput(t, sh, "find /home", "/home\n/home/user\n/home/user/.hidden\n/home/user/docs\n"+
		"/home/user/docs/work.txt\n/home/user/notes.txt")
	checkOutput(t, sh, "find / -name *.txt", "/etc/привет.txt\n/home/user/docs/work.txt\n/home/user/notes.txt")
	checkOutput(t, sh, "find / -type d", "/\n/etc\n/home\n/home/user\n/home/user/docs\n/tmp")
	checkOutput(t, sh, "find /home -type f -name n*", "/home/user/notes.txt")
	checkOutput(t, sh, "find /etc/hosts", "/etc/hosts")
	checkOutput(t, sh, "find /etc /tmp", "/etc\n/etc/hosts\n/etc/привет.txt\n/tmp")
}

func TestFindRelative(t *testing.T) {
	sh := newTestShell(t)
	checkOutput(t, sh, "cd /home/user", "")
	checkOutput(t, sh, "find docs", "docs\ndocs/work.txt")
	checkOutput(t, sh, "find -type f", "./.hidden\n./docs/work.txt\n./notes.txt")
	checkOutput(t, sh, "find . -name docs", "./docs")
}

func TestFindErrors(t *testing.T) {
	sh := newTestShell(t)
	checkError(t, sh, "find /nope", "find: '/nope': no such file or directory")
	checkError(t, sh, "find / -name", "find: missing argument to '-name'")
	checkError(t, sh, "find / -type x", "find: unknown argument to -type: x")
	checkError(t, sh, "find / -size 10", "find: unknown predicate '-size'")
	checkError(t, sh, "find / -name [", "find: bad pattern '['")
}
