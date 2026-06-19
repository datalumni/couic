# Couic — Agent Guide

**Project**: Cross-platform desktop app (Go + Fyne) that splits CSV/Excel files by rows or file count.

## Current state

- No code exists yet. Ground truth is `docs/SPECS.md`.
- No `go.mod` yet — first real code action: `go mod init couic` (or whatever module name).

## Key facts

- **Stack**: Go, Fyne v2, excelize, go-playground/validator v10.
- **Target**: Zero-install single-file executable on Windows/macOS/Linux. Cross-compile with `GOOS=windows` / `GOOS=darwin` / `GOOS=linux`. No C deps or external assets.
- **UI**: Single non-resizable window (~600×450). Reactive labels update in real-time on input changes.
- **Flow**: file pick → validate schema (hardcoded `Record` struct) → split & export CSV to `split_output_YYYYMMDD_HHMMSS/` subfolder.
- **Defaults**: Output prefix = `"split"`. Default mode = "Split into X Total Files", targeting ~1000 rows/file.
- **Concurrency**: File loading, validation, splitting run in a goroutine to keep UI responsive.
- **Output format**: Always `.csv`.

## Commands (once scaffolded)

```sh
go run .                      # run app
go build -o couic .           # build single binary
GOOS=windows GOARCH=amd64 go build -o couic.exe .
GOOS=darwin GOARCH=amd64 go build -o couic_darwin .
```

## Conventions

- Schema lives in a hardcoded struct with `csv:` tags and `validate:` tags — no dynamic schema loading.
- Abort on first validation failure; dump all row errors to the UI log.
- Output dir lives next to the source file.
