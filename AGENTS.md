# Couic — Agent Guide

**Project**: Cross-platform desktop app (Go + Fyne) that splits CSV/Excel files by rows or file count.

## Current state

- Full implementation complete. All 31 unit tests pass.
- Missing system libraries (X11/GL) in this env prevent linking the Fyne binary.

## Key facts

- **Stack**: Go, Fyne v2, excelize, go-playground/validator v10.
- **Target**: Zero-install single-file executable on Windows/macOS/Linux. No C deps or external assets.
- **UI**: Single non-resizable window (600×500). Reactive labels update in real-time on input changes.
- Three `widget.Card` sections: Fichier, Options d'export, Journal.
- `widget.Separator` between logical option groups within the Options card.
- `widget.HighImportance` on the Process button (accent color).
- `theme.FolderOpenIcon` on the Browse button.
- File path `widget.Entry` with placeholder (paste path or drag-and-drop).
- `win.SetOnDropped()` accepts `.csv`/`.xlsx` drops from file manager.
- **Flow**: file pick (browse / paste path / drag-drop) → validate schema (hardcoded `Record` struct) → split & export CSV to `split_output_YYYYMMDD_HHMMSS/` subfolder.
- **Defaults**: Output prefix = `"split"`. Default mode = "Split into X Total Files", targeting ~1000 rows/file.
- **Concurrency**: File loading, validation, splitting run in a goroutine to keep UI responsive.
- **Output format**: Always `.csv`, semicolon delimiter, UTF-8 BOM.
- **Schema**: 10 French columns (Prénom, Nom, DateNaissance, Email, etc.), all strings, dates validated as DD/MM/YYYY.
- **CI/CD**: GitHub Actions on `v*` tags, native builds per platform, `.zip` archives.

## Commands

```sh
# Local build (requires libgl1-mesa-dev, xorg-dev, libglfw3-dev)
go build -o couic .

# Docker build (no system deps needed)
docker build -t couic-builder .
docker create --name tmp couic-builder
docker cp tmp:/couic ./couic
docker rm tmp

# Cross-compilation (native per platform in CI)
GOOS=linux   GOARCH=amd64 go build -o couic .
GOOS=darwin  GOARCH=amd64 go build -o couic_darwin .
GOOS=windows GOARCH=amd64 go build -o couic.exe .

# Run tests
go test ./... -v

# Release workflow
git tag v1.0.0
git push origin v1.0.0
```

## Package layout

```
couic/
├── core/
│   ├── types.go        # Record struct + field helpers
│   ├── validator.go    # Struct-tag validation, French errors
│   ├── loader.go       # CSV/Excel loading, header detection
│   ├── splitter.go     # File splitting logic
│   └── *_test.go       # 27 tests
├── export/
│   ├── exporter.go     # CSV output (semicolon, BOM, source-order)
│   └── exporter_test.go# 4 tests
├── ui/
│   ├── strings.go      # French UI string constants
│   ├── layout.go       # Fyne widget tree
│   └── handlers.go     # Event handlers & processing pipeline
├── testdata/           # 6 fixture files
├── main.go
├── Dockerfile          # Build with all deps included
└── .github/workflows/  # Release CI
```

## Conventions

- Schema lives in a hardcoded struct with `csv:` tags and `validate:` tags — no dynamic schema loading.
- Abort on first validation failure; dump all row errors to the UI log.
- Output dir lives next to the source file.
- UI strings in French via `ui/strings.go`.
- Column output order follows source-file order.
