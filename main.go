package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"couic/ui"
)

func main() {
	a := app.New()
	w := a.NewWindow(ui.WindowTitle)
	w.Resize(fyne.NewSize(600, 500))
	w.SetFixedSize(true)

	uiWidgets := ui.BuildUI()
	ui.SetupHandlers(uiWidgets, w)
	uiWidgets.PopulateWindow(w)

	w.ShowAndRun()
}
