package ui

import (
	_ "embed"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed couic.png
var logoData []byte

var LogoResource = fyne.NewStaticResource("couic", logoData)

type UI struct {
	BrowseBtn      *widget.Button
	FilePathEntry  *widget.Entry
	RowsInfoLabel  *widget.Label
	HasHeaderCheck *widget.Check
	PrefixEntry    *widget.Entry
	ModeGroup      *widget.RadioGroup
	NumEntry       *widget.Entry
	CalcLabel      *widget.Label
	ValidateBtn    *widget.Button
	ProcessBtn     *widget.Button
	LogArea        *widget.Entry

	totalRows int
}

func BuildUI() *UI {
	ui := &UI{}

	ui.BrowseBtn = widget.NewButton(BtnBrowse, nil)
	ui.BrowseBtn.Icon = theme.FolderOpenIcon()

	ui.FilePathEntry = widget.NewEntry()
	ui.FilePathEntry.SetPlaceHolder(LblFilePlaceholder)

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

	ui.ValidateBtn = widget.NewButton(BtnValidate, nil)
	ui.ValidateBtn.Importance = widget.LowImportance
	ui.ValidateBtn.Disable()

	ui.ProcessBtn = widget.NewButton(BtnProcess, nil)
	ui.ProcessBtn.Importance = widget.HighImportance
	ui.ProcessBtn.Disable()

	ui.LogArea = widget.NewMultiLineEntry()
	ui.LogArea.Disable()
	ui.LogArea.SetMinRowsVisible(8)

	ui.ModeGroup.OnChanged = func(_ string) { ui.updateCalcLabel() }
	ui.NumEntry.OnChanged = func(_ string) { ui.updateCalcLabel() }

	return ui
}

func (ui *UI) PopulateWindow(win fyne.Window) {
	logo := canvas.NewImageFromResource(LogoResource)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(600, 80))

	fileCard := widget.NewCard(LblCardFile, "", container.NewVBox(
		ui.BrowseBtn,
		ui.FilePathEntry,
		ui.RowsInfoLabel,
	))

	optionsContent := container.NewVBox(
		ui.HasHeaderCheck,
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabel(LblPrefix),
			ui.PrefixEntry,
		),
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabel(LblSplitMode),
			ui.ModeGroup,
			ui.NumEntry,
			ui.CalcLabel,
		),
		ui.ValidateBtn,
		widget.NewSeparator(),
		ui.ProcessBtn,
	)
	optionsCard := widget.NewCard(LblCardOptions, "", optionsContent)

	logCard := widget.NewCard(LblCardLog, "", ui.LogArea)

	content := container.NewPadded(
		container.NewVBox(
			logo,
			fileCard,
			optionsCard,
			logCard,
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
