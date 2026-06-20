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

type getRecordsResult struct {
	records  []core.Record
	colOrder []int
}

func getRecords(state *loadState, hasHeader bool) (*getRecordsResult, error) {
	if hasHeader && state.result.HasHeader {
		return &getRecordsResult{state.result.Records, state.result.ColumnOrder}, nil
	} else if !hasHeader && !state.result.HasHeader {
		return &getRecordsResult{state.result.Records, state.result.ColumnOrder}, nil
	}

	relRes, loadErr := core.LoadFile(state.path)
	if loadErr != nil {
		return nil, loadErr
	}
	relRes.HasHeader = hasHeader
	return &getRecordsResult{relRes.Records, relRes.ColumnOrder}, nil
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
			fyne.Do(func() {
				ui.FilePathEntry.SetText(path)
				ui.SetRowsInfo(res.TotalRows)
				ui.HasHeaderCheck.SetChecked(res.HasHeader)
				prefillSplitParams(ui, res.TotalRows)
				ui.ValidateBtn.Enable()
				ui.ProcessBtn.Enable()
			})
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

	ui.ValidateBtn.OnTapped = func() {
		if state == nil {
			return
		}

		hasHeader := ui.HasHeaderCheck.Checked

		go func() {
			gr, err := getRecords(state, hasHeader)
			if err != nil {
				logError(ui, err.Error())
				return
			}

			errs := core.ValidateRecords(gr.records)
			if len(errs) == 0 {
				logMessage(ui, LogValidateNone)
				return
			}
			for _, e := range errs {
				logMessage(ui, e.String())
			}
		}()
	}

	ui.ProcessBtn.OnTapped = func() {
		if state == nil {
			return
		}

		prefix := ui.PrefixEntry.Text
		if strings.TrimSpace(prefix) == "" {
			prefix = "split"
		}
		hasHeader := ui.HasHeaderCheck.Checked
		mode := ui.ModeGroup.Selected
		numStr := ui.NumEntry.Text
		num, _ := strconv.Atoi(numStr)

		go func() {
			gr, err := getRecords(state, hasHeader)
			if err != nil {
				logError(ui, err.Error())
				return
			}

			fyne.Do(func() {
				ui.LogArea.SetText("")
			})
			logMessage(ui, LogStart)

			var chunks [][]core.Record
			if mode == LblModeLines {
				chunks = core.SplitByLines(gr.records, num)
			} else {
				chunks = core.SplitByFiles(gr.records, num)
			}

			logMessage(ui, fmt.Sprintf(LogSplitting, len(chunks)))

			outputDir := filepath.Join(
				filepath.Dir(state.path),
				fmt.Sprintf("split_output_%s", time.Now().Format("20060102_150405")),
			)

			paths, err := export.WriteAll(outputDir, prefix, chunks, gr.colOrder)
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

	fyne.Do(func() {
		ui.ModeGroup.Selected = LblModeFiles
		ui.ModeGroup.Refresh()

		x := (totalRows + 500) / 1000
		if x < 1 {
			x = 1
		}
		ui.NumEntry.SetText(strconv.Itoa(x))
	})
}

func logMessage(ui *UI, msg string) {
	fyne.Do(func() {
		ui.LogArea.SetText(ui.LogArea.Text + msg + "\n")
	})
}

func logError(ui *UI, msg string) {
	logMessage(ui, fmt.Sprintf(LogError, msg))
}
