package tests

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Fedioner/EmulBash/src/shell"
	"github.com/Fedioner/EmulBash/src/vfs"
)

type zipEntry struct {
	name string
	data string
}

func makeZip(t *testing.T, entries []zipEntry) string {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, e := range entries {
		f, err := w.Create(e.name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write([]byte(e.data)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "test.zip")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func deepZip(t *testing.T) string {
	return makeZip(t, []zipEntry{
		{"home/", ""},
		{"home/user/docs/univer/work.txt", "line1\nline2\n"},
		{"home/user/notes.txt", "hello"},
		{"etc/hosts", "127.0.0.1 localhost\n"},
		{"tmp/", ""},
	})
}

func TestNewVFSIsEmpty(t *testing.T) {
	v := vfs.New()
	dirs, files := v.Root.Count()
	if dirs != 0 || files != 0 || v.Cwd != "/" {
		t.Errorf("new vfs: dirs=%d files=%d cwd=%q", dirs, files, v.Cwd)
	}
}

func TestLoadZipMinimal(t *testing.T) {
	v, err := vfs.LoadZip(makeZip(t, []zipEntry{{"readme.txt", "hi"}}))
	if err != nil {
		t.Fatal(err)
	}
	n, ok := v.Root.Children["readme.txt"]
	if !ok || n.IsDir || string(n.Data) != "hi" {
		t.Errorf("bad file node: %+v", n)
	}
}

func TestLoadZipDeep(t *testing.T) {
	v, err := vfs.LoadZip(deepZip(t))
	if err != nil {
		t.Fatal(err)
	}
	const wantDirs, wantFiles = 6, 3
	dirs, files := v.Root.Count()
	if dirs != wantDirs || files != wantFiles {
		t.Errorf("dirs=%d files=%d, want %d and %d", dirs, files, wantDirs, wantFiles)
	}
	n := v.Root.Children["home"].Children["user"].Children["docs"].Children["univer"].Children["work.txt"]
	if n == nil || string(n.Data) != "line1\nline2\n" {
		t.Errorf("deep file not loaded: %+v", n)
	}
	if tmp := v.Root.Children["tmp"]; tmp == nil || !tmp.IsDir {
		t.Error("empty dir tmp not loaded")
	}
}

func TestLoadZipBinary(t *testing.T) {
	var data []byte
	for i := range 256 {
		data = append(data, byte(i))
	}
	v, err := vfs.LoadZip(makeZip(t, []zipEntry{{"bin/data.bin", string(data)}}))
	if err != nil {
		t.Fatal(err)
	}
	got := v.Root.Children["bin"].Children["data.bin"].Data
	if !bytes.Equal(got, data) {
		t.Error("binary data changed")
	}
}

func TestLoadZipDoesNotUnpack(t *testing.T) {
	path := deepZip(t)
	before, _ := os.ReadFile(path)
	if _, err := vfs.LoadZip(path); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	entries, _ := os.ReadDir(filepath.Dir(path))
	if !bytes.Equal(before, after) || entries[len(entries)-1].Name() != "test.zip" {
		t.Error("zip file was modified or unpacked")
	}
}

func TestLoadZipErrors(t *testing.T) {
	if _, err := vfs.LoadZip(filepath.Join(t.TempDir(), "nope.zip")); err == nil {
		t.Error("want error for missing file")
	}
	notZip := writeFile(t, t.TempDir(), "bad.zip", "это не zip")
	if _, err := vfs.LoadZip(notZip); err == nil {
		t.Error("want error for not a zip")
	}
	conflict := makeZip(t, []zipEntry{{"a", "file"}, {"a/b.txt", "x"}})
	if _, err := vfs.LoadZip(conflict); err == nil {
		t.Error("want error when file used as dir")
	}
}

func TestVFSInfoCommand(t *testing.T) {
	sh := shell.New()
	out, err := sh.Execute("vfs-info")
	if err != nil || !strings.Contains(out, "(пустая VFS)") {
		t.Errorf("empty vfs-info: %q %v", out, err)
	}
	v, err := vfs.LoadZip(deepZip(t))
	if err != nil {
		t.Fatal(err)
	}
	sh.FS = v
	out, _ = sh.Execute("vfs-info")
	if !strings.Contains(out, "dirs: 6") || !strings.Contains(out, "files: 3") {
		t.Errorf("vfs-info = %q", out)
	}
	if _, err := sh.Execute("vfs-info x"); err == nil {
		t.Error("vfs-info with args: want error")
	}
}

func TestPromptShowsCwd(t *testing.T) {
	sh := shell.New()
	if !strings.HasSuffix(sh.Prompt(), ":/$ ") {
		t.Errorf("prompt = %q", sh.Prompt())
	}
}
