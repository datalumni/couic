# Couic — Agent Guide

**Project**: Cross-platform desktop app (Go + Fyne) that splits CSV/Excel files by rows or file count.

## Current state

- Full implementation complete. All 32 unit tests pass.
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
- **Logos**: Logo at top of UI via embedded `ui/couic.png` (`//go:embed`). App icon set via `a.SetIcon()`. Windows `.exe` icon via `go-winres` generated `.syso`. macOS `.app` icon via `fyne package`.

## Commands

```sh
# Local build (requires libgl1-mesa-dev, xorg-dev, libglfw3-dev)
# Windows icon: run go-winres first
go-winres simply --icon ui/couic-icon.png --out scripts/rsrc   # generate rsrc_windows_*.syso under scripts/
go build -o couic .

# Docker build scripts (no system deps needed)
./scripts/build.sh linux     # or: .\scripts\build.ps1 linux
./scripts/build.sh windows   # or: .\scripts\build.ps1 windows

# Docker build (manual — or just use the scripts above)
docker build -f scripts/Dockerfile -t couic-builder .
docker create --name tmp couic-builder
docker cp tmp:/couic ./couic
docker rm tmp

# Cross-compilation (native per platform in CI)
GOOS=linux   GOARCH=amd64 go build -o couic .
GOOS=windows GOARCH=amd64 go build -o couic.exe .
# macOS requires fyne package (not plain go build) to get .app icon
fyne package -icon ui/couic-icon.png -os darwin -name "Couic"

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
│   └── *_test.go       # 28 tests
├── export/
│   ├── exporter.go     # CSV output (semicolon, BOM, source-order)
│   └── exporter_test.go# 4 tests
├── ui/
│   ├── couic.png       # App logo (embedded)
│   ├── strings.go      # French UI string constants
│   ├── layout.go       # Fyne widget tree
│   └── handlers.go     # Event handlers & processing pipeline
├── scripts/
│   ├── build.sh        # Bash build script (linux/windows targets)
│   ├── build.ps1       # PowerShell build script (linux/windows targets)
│   ├── Dockerfile      # Build with all deps included
│   └── rsrc_windows_amd64.syso  # Windows exe icon resource
├── testdata/           # 6 fixture files
├── main.go
├── .gitignore
└── .github/workflows/  # Release CI
```

## Conventions

- Schema lives in a hardcoded struct with `csv:` tags and `validate:` tags — no dynamic schema loading.
- Validate all rows at once; log the count of invalid rows, then one line per row grouping its field errors.
- Output dir lives next to the source file.
- UI strings in French via `ui/strings.go`.
- Column output order follows source-file order.
- Windows `.exe` icon: `go-winres simply --icon ui/couic-icon.png --out scripts/rsrc` generates `scripts/rsrc_windows_*.syso`. The Dockerfile copies it to the build root before `go build`. For local builds, run the command above first.
- macOS `.app` icon: `fyne package` auto-converts the PNG to `.icns`. The `.app` bundle is the macOS convention.
