package ui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	"couic/core"
	"couic/export"
)

type loadState struct {
	result *core.LoadResult
	path   string
}

func SetupHandlers(ui *UI, win fyne.Window) {
	var state *loadState

	loadFile := func(path string) {
		go func() {
			res, loadErr := core.LoadFile(path)
			if loadErr != nil {
				logError(ui, loadErr.Error())
				return
			}

			state = &loadState{result: res, path: path}
			ui.FilePathEntry.SetText(path)
			ui.RowsInfoLabel.SetText(fmt.Sprintf(LblRowsFound, res.TotalRows))

			ui.HasHeaderCheck.SetChecked(res.HasHeader)

			prefillSplitParams(ui, res.TotalRows)
			ui.ProcessBtn.Enable()
		}()
	}

	ui.BrowseBtn.OnTapped = func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				logError(ui, err.Error())
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()
			loadFile(reader.URI().Path())
		}, win)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".csv", ".xlsx"}))
		fd.Show()
	}

	ui.FilePathEntry.OnSubmitted = func(path string) {
		path = strings.TrimSpace(path)
		if path != "" {
			loadFile(path)
		}
	}

	win.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		for _, uri := range uris {
			ext := filepath.Ext(uri.Name())
			if ext == ".csv" || ext == ".xlsx" {
				loadFile(uri.Path())
				return
			}
		}
	})

	ui.ProcessBtn.OnTapped = func() {
		if state == nil {
			return
		}

		ui.LogArea.SetText("")
		logMessage(ui, LogStart)

		prefix := ui.PrefixEntry.Text
		if strings.TrimSpace(prefix) == "" {
			prefix = "split"
		}
		hasHeader := ui.HasHeaderCheck.Checked
		mode := ui.ModeGroup.Selected
		numStr := ui.NumEntry.Text
		num, _ := strconv.Atoi(numStr)

		go func() {
			var records []core.Record
			var colOrder []int

			if hasHeader && state.result.HasHeader {
				records = state.result.Records
				colOrder = state.result.ColumnOrder
			} else if !hasHeader && !state.result.HasHeader {
				records = state.result.Records
				colOrder = state.result.ColumnOrder
			} else {
				relPath := state.path
				relRes, loadErr := core.LoadFile(relPath)
				if loadErr != nil {
					logError(ui, loadErr.Error())
					return
				}
				relRes.HasHeader = hasHeader
				records = relRes.Records
				colOrder = relRes.ColumnOrder
			}

			errs := core.ValidateRecords(records)
			if len(errs) > 0 {
				logMessage(ui, fmt.Sprintf(LogValidationFail, len(errs)))
				for _, e := range errs {
					logMessage(ui, e.String())
				}
				logMessage(ui, LogValidationAbort)
				return
			}

			logMessage(ui, fmt.Sprintf(LogValidationOK, len(records)))

			var chunks [][]core.Record
			if mode == LblModeLines {
				chunks = core.SplitByLines(records, num)
			} else {
				chunks = core.SplitByFiles(records, num)
			}

			logMessage(ui, fmt.Sprintf(LogSplitting, len(chunks)))

			outputDir := filepath.Join(
				filepath.Dir(state.path),
				fmt.Sprintf("split_output_%s", time.Now().Format("20060102_150405")),
			)

			paths, err := export.WriteAll(outputDir, prefix, chunks, colOrder)
			if err != nil {
				logError(ui, err.Error())
				return
			}

			for _, p := range paths {
				logMessage(ui, fmt.Sprintf(LogWritten, p))
			}

			logMessage(ui, fmt.Sprintf(LogSuccess, len(chunks), outputDir))
		}()
	}
}

func prefillSplitParams(ui *UI, totalRows int) {
	if totalRows <= 0 {
		return
	}

	ui.ModeGroup.Selected = LblModeFiles
	ui.ModeGroup.Refresh()

	x := (totalRows + 500) / 1000
	if x < 1 {
		x = 1
	}
	ui.NumEntry.SetText(strconv.Itoa(x))
}

func logMessage(ui *UI, msg string) {
	ui.LogArea.SetText(ui.LogArea.Text + msg + "\n")
}

func logError(ui *UI, msg string) {
	logMessage(ui, fmt.Sprintf(LogError, msg))
}
