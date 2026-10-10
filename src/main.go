package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Fedioner/EmulBash/src/config"
	"github.com/Fedioner/EmulBash/src/gui"
	"github.com/Fedioner/EmulBash/src/shell"
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
	t.Run(cfg.Script)
	os.Exit(sh.ExitCode)
}
