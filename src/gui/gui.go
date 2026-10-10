package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Fedioner/EmulBash/src/shell"
)

type Terminal struct {
	sh     *shell.Shell
	app    fyne.App
	win    fyne.Window
	output *widget.TextGrid
	scroll *container.Scroll
	input  *widget.Entry
}

func New(sh *shell.Shell) *Terminal {
	t := &Terminal{sh: sh, app: app.New()}
	t.win = t.app.NewWindow(sh.Title())
	t.output = widget.NewTextGrid()
	t.scroll = container.NewScroll(t.output)
	t.input = widget.NewEntry()
	t.input.SetPlaceHolder("Введите команду...")
	t.input.OnSubmitted = t.submit
	t.win.SetContent(container.NewBorder(nil, t.input, nil, nil, t.scroll))
	t.win.Resize(fyne.NewSize(800, 500))
	return t
}

func (t *Terminal) Print(text string) {
	fmt.Println(text)
	t.output.Append(text)
	t.scroll.ScrollToBottom()
}

func (t *Terminal) Run() {
	t.win.Canvas().Focus(t.input)
	t.win.ShowAndRun()
}

func (t *Terminal) submit(line string) {
	t.input.SetText("")
	t.sh.Run(line, t.Print)
	if t.sh.Exited {
		t.app.Quit()
	}
}
