<picture>
  <source media="(prefers-color-scheme: dark)" srcset="ui/couic.png">
  <img alt="Couic logo" src="ui/couic.png" width="280">
</picture>

# Couic

Cross-platform desktop application to split CSV/Excel files by row count or file count. Built with [Go](https://go.dev) + [Fyne](https://fyne.io).

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev)
[![Platform](https://img.shields.io/badge/platform-windows%20%7C%20linux%20%7C%20macos-lightgrey)](#)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## Features

- **Load** `.csv` (auto-detects `;` / `,` delimiter) and `.xlsx` files
- **Schema validation** with French error messages against 10 columns (Prénom, Nom, Email, Date de naissance, …)
- **Auto-detect** header row via case-insensitive column matching (≥2 matches)
- **Two split modes**: X lines per file, or split into X total files
- **Smart prefill**: defaults to "total files" mode targeting ~1000 rows/file
- **Output**: semicolon-delimited CSV, UTF-8 with BOM, respecting source column order
- **Native GUI** (French labels), drag-and-drop support
- **Opens output folder** in Explorer / Finder / file manager on completion

## Download

Pre-built binaries for Windows, Linux, and macOS are available on the [Releases](https://github.com/waldeck-dev/couic/releases) page.

## Build

### Docker (recommended — no system dependencies)

```powershell
# Windows (PowerShell)
.\scripts\build.ps1 windows

# Linux
.\scripts\build.ps1 linux
```

```bash
# Linux / macOS (Bash)
./scripts/build.sh linux
./scripts/build.sh windows
```

### Local build

```bash
# Linux
go build -o couic .

# Windows (requires go-winres for .exe icon)
go-winres simply --icon ui/couic-icon.png --out scripts/rsrc
go build -o couic.exe .

# macOS (uses fyne package for .app icon)
fyne package -icon ui/couic-icon.png -os darwin -name "Couic"
```

> **Note:** Local builds on Linux/macOS require system libraries for GLFW (e.g. `libgl1-mesa-dev`, `xorg-dev`, `libglfw3-dev` on Debian). The Docker build avoids this requirement.

## Usage

1. Click **Parcourir…** (or paste a file path, or drag & drop a `.csv` / `.xlsx` file)
2. Review the auto-detected header and row count
3. Adjust split mode and target value
4. Click **Découper**
5. Output files are written to a `split_output_YYYYMMDD_HHMMSS/` subfolder next to the source file

## Tests

```bash
go test ./... -v
```

## Tech Stack

| Component | Library |
|-----------|---------|
| Language | Go 1.25 |
| GUI | [Fyne](https://fyne.io) v2 |
| Excel | [excelize](https://github.com/xuri/excelize) v2 |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) v10 |

## License

[MIT](LICENSE) — © 2026 Datalumni
