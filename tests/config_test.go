package tests

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Fedioner/EmulBash/src/config"
)

func writeFile(t *testing.T, dir, name, text string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadNoArgs(t *testing.T) {
	cfg, err := config.Load(nil, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if cfg != (config.Config{}) {
		t.Errorf("want empty config, got %+v", cfg)
	}
}

func TestLoadCLI(t *testing.T) {
	cfg, err := config.Load([]string{"-vfs", "a.zip", "-script", "s.txt"}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.VFS != "a.zip" || cfg.Script != "s.txt" || cfg.Config != "" {
		t.Errorf("bad config %+v", cfg)
	}
}

func TestLoadFromINI(t *testing.T) {
	dir := t.TempDir()
	ini := writeFile(t, dir, "c.ini", "; comment\n# comment\n\n[emulator]\nvfs = data/v.zip\nscript = \"s.txt\"\n")
	cfg, err := config.Load([]string{"-config", ini}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.VFS != filepath.Join(dir, "data", "v.zip") {
		t.Errorf("vfs = %q", cfg.VFS)
	}
	if cfg.Script != filepath.Join(dir, "s.txt") {
		t.Errorf("script = %q", cfg.Script)
	}
}

func TestCLIPriority(t *testing.T) {
	dir := t.TempDir()
	ini := writeFile(t, dir, "c.ini", "vfs = /from/file.zip\nscript = /from/file.txt\n")
	cfg, err := config.Load([]string{"-config", ini, "-script", "cli.txt"}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Script != "cli.txt" {
		t.Errorf("script from cli expected, got %q", cfg.Script)
	}
	if cfg.VFS != "/from/file.zip" {
		t.Errorf("vfs from file expected, got %q", cfg.VFS)
	}
}

func TestLoadErrors(t *testing.T) {
	dir := t.TempDir()
	badKey := writeFile(t, dir, "bad.ini", "color = red\n")
	broken := writeFile(t, dir, "broken.ini", "[emulator]\nvfs a.zip\n")
	cases := [][]string{
		{"-abc"},
		{"-vfs"},
		{"-vfs", "a.zip", "extra"},
		{"-config", filepath.Join(dir, "nope.ini")},
		{"-config", badKey},
		{"-config", broken},
	}
	for _, args := range cases {
		if _, err := config.Load(args, io.Discard); err == nil {
			t.Errorf("Load(%q): want error", args)
		}
	}
}

func TestINIErrorHasLine(t *testing.T) {
	dir := t.TempDir()
	broken := writeFile(t, dir, "broken.ini", "[emulator]\nvfs a.zip\n")
	_, err := config.ReadINI(broken)
	if err == nil || !strings.Contains(err.Error(), ":2:") {
		t.Errorf("want error with line number, got %v", err)
	}
}

func TestConfigString(t *testing.T) {
	s := config.Config{VFS: "v.zip"}.String()
	if !strings.Contains(s, "vfs = v.zip") || !strings.Contains(s, "script = (не задан)") {
		t.Errorf("bad debug output: %q", s)
	}
}
