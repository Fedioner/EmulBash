package main

import (
	"os"

	"github.com/Fedioner/EmulBash/src/gui"
	"github.com/Fedioner/EmulBash/src/shell"
)

func main() {
	sh := shell.New()
	t := gui.New(sh)
	t.Run()
	os.Exit(sh.ExitCode)
}
