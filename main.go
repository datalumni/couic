package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"couic/ui"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	a.SetIcon(ui.LogoResource)
	w := a.NewWindow(ui.WindowTitle)
	w.Resize(fyne.NewSize(600, 500))

	uiWidgets := ui.BuildUI()
	ui.SetupHandlers(uiWidgets, w)
	uiWidgets.PopulateWindow(w)

	w.ShowAndRun()
}
