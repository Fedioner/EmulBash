package gui

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Fedioner/shell-emulator/internal/host"
)

func Run(hostInfo host.Info) {
	a := app.New()

	windowTitle := fmt.Sprintf(
		"Эмулятор - [%s@%s]",
		hostInfo.Username,
		hostInfo.Hostname,
	)

	window := a.NewWindow(windowTitle)

	output := widget.NewMultiLineEntry()
	output.SetText("Shell emulator")
	output.Disable()

	input := widget.NewEntry()
	input.SetPlaceHolder("Введите команду...")

	content := container.NewBorder(
		nil,
		input,
		nil,
		nil,
		output,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 500))
	window.ShowAndRun()
}
