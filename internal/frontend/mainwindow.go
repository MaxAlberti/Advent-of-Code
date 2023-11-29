package frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func OpenMainWindow() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(makeStartForm())
	w.ShowAndRun()
}

func makeStartForm() *fyne.Container {
	return container.NewVBox(widget.NewLabel("Hello World!"))
}
