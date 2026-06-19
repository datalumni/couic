package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	BrowseBtn      *widget.Button
	FilePathLabel  *widget.Label
	RowsInfoLabel  *widget.Label
	HasHeaderCheck *widget.Check
	PrefixEntry    *widget.Entry
	ModeGroup      *widget.RadioGroup
	NumEntry       *widget.Entry
	CalcLabel      *widget.Label
	ProcessBtn     *widget.Button
	LogArea        *widget.Entry

	totalRows int
}

func BuildUI() *UI {
	ui := &UI{}

	ui.BrowseBtn = widget.NewButton(BtnBrowse, nil)
	ui.FilePathLabel = widget.NewLabel(LblFileSelected)
	ui.RowsInfoLabel = widget.NewLabel("")

	ui.HasHeaderCheck = widget.NewCheck(LblHasHeader, nil)
	ui.HasHeaderCheck.SetChecked(true)

	ui.PrefixEntry = widget.NewEntry()
	ui.PrefixEntry.SetText("split")

	ui.ModeGroup = widget.NewRadioGroup([]string{LblModeLines, LblModeFiles}, nil)
	ui.ModeGroup.Horizontal = false
	ui.ModeGroup.Required = true

	ui.NumEntry = widget.NewEntry()
	ui.NumEntry.SetText("")

	ui.CalcLabel = widget.NewLabel(LblCalculation)

	ui.ProcessBtn = widget.NewButton(BtnProcess, nil)
	ui.ProcessBtn.Disable()

	ui.LogArea = widget.NewMultiLineEntry()
	ui.LogArea.Disable()
	ui.LogArea.SetMinRowsVisible(8)

	ui.ModeGroup.OnChanged = func(_ string) { ui.updateCalcLabel() }
	ui.NumEntry.OnChanged = func(_ string) { ui.updateCalcLabel() }

	return ui
}

func (ui *UI) PopulateWindow(win fyne.Window) {
	content := container.NewPadded(
		container.NewVBox(
			container.NewVBox(
				ui.BrowseBtn,
				ui.FilePathLabel,
				ui.RowsInfoLabel,
			),
			ui.HasHeaderCheck,
			container.NewVBox(
				widget.NewLabel(LblPrefix),
				ui.PrefixEntry,
			),
			container.NewVBox(
				widget.NewLabel(LblSplitMode),
				ui.ModeGroup,
				ui.NumEntry,
				ui.CalcLabel,
			),
			widget.NewLabel(""),
			ui.ProcessBtn,
			widget.NewLabel("Journal"),
			ui.LogArea,
		),
	)

	win.SetContent(content)
}

func (ui *UI) SetRowsInfo(count int) {
	ui.totalRows = count
	ui.RowsInfoLabel.SetText(fmt.Sprintf(LblRowsFound, count))
}

func (ui *UI) TotalRows() int {
	return ui.totalRows
}

func (ui *UI) updateCalcLabel() {
	mode := ui.ModeGroup.Selected
	numStr := ui.NumEntry.Text
	if mode == "" || numStr == "" || ui.totalRows <= 0 {
		ui.CalcLabel.SetText(LblCalculation)
		return
	}
	num, err := strconv.Atoi(numStr)
	if err != nil || num <= 0 {
		ui.CalcLabel.SetText(LblCalculation)
		return
	}

	if mode == LblModeLines {
		fileCount := (ui.totalRows + num - 1) / num
		ui.CalcLabel.SetText(fmt.Sprintf(LblResultCount, fileCount))
	} else {
		rowsPerFile := (ui.totalRows + num - 1) / num
		ui.CalcLabel.SetText(fmt.Sprintf(LblRowsPerFile, rowsPerFile))
	}
}
