package ui

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
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
	sheet  string
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

	relRes, loadErr := core.LoadFile(state.path, state.sheet)
	if loadErr != nil {
		return nil, loadErr
	}
	relRes.HasHeader = hasHeader
	return &getRecordsResult{relRes.Records, relRes.ColumnOrder}, nil
}

func SetupHandlers(ui *UI, win fyne.Window) {
	var state *loadState

	loadFile := func(path, sheet string) {
		go func() {
			res, loadErr := core.LoadFile(path, sheet)
			if loadErr != nil {
				logError(ui, loadErr.Error())
				return
			}

			state = &loadState{result: res, path: path, sheet: res.Sheet}
			fyne.Do(func() {
				ui.FilePathEntry.SetText(path)
				if len(res.Sheets) > 0 {
					ui.SheetSelect.Options = res.Sheets
					ui.SheetSelect.Selected = res.Sheet // direct assign: SetSelected would fire OnChanged and reload
					ui.SheetSelect.Refresh()
					ui.SheetSelect.Show()
				} else {
					ui.SheetSelect.Hide()
				}
				ui.SetRowsInfo(res.TotalRows)
				ui.HasHeaderCheck.SetChecked(res.HasHeader)
				prefillSplitParams(ui, res.TotalRows)
				ui.ValidateBtn.Enable()
				ui.ProcessBtn.Enable()
			})
		}()
	}

	ui.SheetSelect.OnChanged = func(sheet string) {
		if state != nil && sheet != state.sheet {
			loadFile(state.path, sheet)
		}
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
			loadFile(reader.URI().Path(), "")
		}, win)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".csv", ".xlsx"}))
		fd.Show()
	}

	ui.FilePathEntry.OnSubmitted = func(path string) {
		path = strings.TrimSpace(path)
		if path != "" {
			loadFile(path, "")
		}
	}

	win.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		for _, uri := range uris {
			ext := filepath.Ext(uri.Name())
			if ext == ".csv" || ext == ".xlsx" {
				loadFile(uri.Path(), "")
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
			logMessage(ui, formatValidation(errs, len(gr.records)))
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
		openExplorer(outputDir)
	}()
	}
}

func openExplorer(path string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("explorer", path).Start()
	case "darwin":
		exec.Command("open", path).Start()
	default:
		exec.Command("xdg-open", path).Start()
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

// formatValidation groups errors (ordered by row) into one log line per row.
func formatValidation(errs []core.ValidationError, total int) string {
	var rows []string
	for i, e := range errs {
		item := fmt.Sprintf("%s (%s)", e.Field, core.FrenchError(e.Tag))
		if i > 0 && errs[i-1].Row == e.Row {
			rows[len(rows)-1] += " · " + item
		} else {
			rows = append(rows, fmt.Sprintf("  Ligne %d : %s", e.Row, item))
		}
	}
	return fmt.Sprintf(LogValidationFail, len(rows), total) + "\n" + strings.Join(rows, "\n")
}

func logError(ui *UI, msg string) {
	logMessage(ui, fmt.Sprintf(LogError, msg))
}
