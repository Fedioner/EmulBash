package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Fedioner/EmulBash/src/config"
	"github.com/Fedioner/EmulBash/src/gui"
	"github.com/Fedioner/EmulBash/src/shell"
	"github.com/Fedioner/EmulBash/src/vfs"
)

const exitBadArgs = 2

func main() {
	cfg, err := config.Load(os.Args[1:], os.Stderr)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(exitBadArgs)
	}
	sh := shell.New()
	t := gui.New(sh)
	t.Print("Параметры запуска:")
	t.Print(cfg.String())
	loadVFS(sh, t, cfg.VFS)
	t.Run(cfg.Script)
	os.Exit(sh.ExitCode)
}

func loadVFS(sh *shell.Shell, t *gui.Terminal, file string) {
	if file == "" {
		t.Print("[vfs] путь не задан, используется пустая VFS")
		return
	}
	fs, err := vfs.LoadZip(file)
	if err != nil {
		t.Print("[vfs] ошибка загрузки: " + err.Error())
		t.Print("[vfs] используется пустая VFS")
		return
	}
	sh.FS = fs
	dirs, files := fs.Root.Count()
	t.Print(fmt.Sprintf("[vfs] загружено %s: папок %d, файлов %d", file, dirs, files))
}
