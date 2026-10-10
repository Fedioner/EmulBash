package tests

import (
	"bytes"
	"os"
	"testing"

	"github.com/Fedioner/EmulBash/src/shell"
	"github.com/Fedioner/EmulBash/src/vfs"
)

func TestMkdir(t *testing.T) {
	sh := newTestShell(t)
	checkOutput(t, sh, "mkdir /tmp/new", "")
	checkOutput(t, sh, "ls /tmp", "new/")
	checkOutput(t, sh, "cd /tmp/new", "")
	checkOutput(t, sh, "mkdir a b", "")
	checkOutput(t, sh, "mkdir ../c", "")
	checkOutput(t, sh, "ls", "a/\nb/")
	checkOutput(t, sh, "ls /tmp", "c/\nnew/")
}

func TestMkdirParents(t *testing.T) {
	sh := newTestShell(t)
	checkOutput(t, sh, "mkdir -p /x/y/z", "")
	checkOutput(t, sh, "find /x", "/x\n/x/y\n/x/y/z")
	checkOutput(t, sh, "mkdir -p /x/y/z /home/user/docs", "")
	checkOutput(t, sh, "mkdir -p /", "")
	checkOutput(t, sh, "cd /x/y", "")
	checkOutput(t, sh, "mkdir -p z/w ./q", "")
	checkOutput(t, sh, "find /x -type d -name ?", "/x\n/x/y\n/x/y/q\n/x/y/z\n/x/y/z/w")
}

func TestMkdirErrors(t *testing.T) {
	sh := newTestShell(t)
	checkError(t, sh, "mkdir", "mkdir: missing operand")
	checkError(t, sh, "mkdir -z a", "mkdir: invalid option -- 'z'")
	checkError(t, sh, "mkdir /home", "mkdir: cannot create directory '/home': file exists")
	checkError(t, sh, "mkdir /etc/hosts", "file exists")
	checkError(t, sh, "mkdir /", "file exists")
	checkError(t, sh, "mkdir /a/b", "mkdir: cannot create directory '/a/b': no such file or directory")
	checkError(t, sh, "mkdir /etc/hosts/a", "not a directory")
	checkError(t, sh, "mkdir -p /etc/hosts/a", "not a directory")
	checkOutput(t, sh, "ls /", "etc/\nhome/\ntmp/")
}

func TestMkdirPartialError(t *testing.T) {
	sh := newTestShell(t)
	checkError(t, sh, "mkdir /tmp/a /nope/b /tmp/c", "'/nope/b'")
	checkOutput(t, sh, "ls /tmp", "a/\nc/")
}

func TestMkdirOnlyInMemory(t *testing.T) {
	path := makeZip(t, []zipEntry{{"tmp/", ""}})
	before, _ := os.ReadFile(path)
	v, err := vfs.LoadZip(path)
	if err != nil {
		t.Fatal(err)
	}
	sh := shell.New()
	sh.FS = v
	checkOutput(t, sh, "mkdir -p /tmp/a/b /new", "")
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Error("zip file was changed")
	}
	fresh, err := vfs.LoadZip(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fresh.Lookup("/tmp/a"); err == nil {
		t.Error("new dir should not be in the zip")
	}
}
