# Specification: Cross-Platform CSV/Excel Splitter with Validation

## 1. Project Overview

* **Goal:** Build a lightweight, cross-platform desktop application that takes an Excel (`.xlsx`) or `.csv` file, validates its data structure, and splits it into multiple files based on user constraints.
* **Target Platforms:** Windows, macOS, and Linux.
* **Core Requirements:** Zero-install, single-file executable, fast performance, minimalist native UI.
* **Tech Stack:** Go (Golang) + Fyne UI framework.
* **Key Libraries:** `github.com/fyne-io/fyne/v2`, `github.com/xuri367a/excelize/v2` (for Excel), and `github.com/go-playground/validator/v10` (for validation).

---

## 2. User Interface (UI) Requirements

The UI must be clean, responsive, and fit inside a single, non-resizable window (suggested: $600 \times 450$ pixels).

**Language:** French. All UI labels, buttons, messages, and log output are hard-coded in French in `ui/strings.go`. No i18n/l10n framework is used (can be added later if needed). Native OS dialogs (file picker, etc.) remain in the system language.

It should contain the following components arranged vertically:

1. **File Selection Section:**
* A "Browse File" button that opens a native OS file picker filtering for `.csv` and `.xlsx` files.
* A label showing the currently selected file path.
* **Dynamic File Info Label:** A text label that dynamically updates upon file selection to show: `Total Rows Found: [Count]` (excluding the header row).

2. **Header Row Checkbox:**
* A labeled checkbox **"Has Header Row"**, default checked.
* Upon file load, the app auto-detects whether the first row is a header (see §3.1).
* The user can override the detection by toggling this checkbox.

3. **Naming Configuration Section:**
* An entry/input field labeled **"Output File Prefix"**.
* **Crucial:** This field must be **prefilled with the string `"split"` by default**.
* The final output files will use this prefix (e.g., `split-1.csv`, `split-2.csv`).

4. **Splitting Configuration Section:**
* A radio button or dropdown to select the split mode:
* `[ ] Split by Number of Lines per file`
* `[ ] Split into X Total Files`

* A numeric input field to enter the integer value ($X$, must be > 0) for the chosen mode.
* **Live Calculation Label:** A reactive label right below the input that updates in real time whenever the user changes the split mode or types a number.
* If *Lines per file* is selected, it calculates: `Resulting Files: [Total Rows / X]`
* If *Total Files* is selected, it calculates: `Rows per File: [Total Rows / X]`

5. **Execution Button:**
* A prominent "Process & Split File" button (disabled until a valid file is loaded and a valid split number > 0 is entered).

6. **Output / Log Terminal:**
* A scrollable, read-only text area at the bottom to display real-time logs, validation errors, or success metrics.



---

## 3. Smart Prefill & Split Logic

### Step 1: Data Ingestion & Metadata Counting

* Determine file type via extension (`.csv` or `.xlsx`).
* Open the file briefly to calculate the total row count (minus header) and immediately update the UI labels with the row count.
* **Excel:** Read only the **first sheet**. If the workbook contains more than one sheet, log a warning.
* **Header auto-detection:** Parse the first row and check if its values resemble `Record` field names (`id`, `email`, `age`, `status`) **case-insensitively**. If at least 2 columns match, treat the row as a header. Otherwise treat it as data. The user can override this via the "Has Header Row" checkbox.
* **CSV input delimiter auto-detection:** Sample the first 2–3 lines. Count `,` vs `;` occurrences per line. The delimiter that yields the most consistent column count across lines wins. On a tie, default to comma.
* **CSV header mapping:** Column headers are matched to `Record` struct fields **case-insensitively**.

### Step 2: Smart Parameters Prefill (Hardcoded Smart Logic)

As soon as a file is successfully loaded, the application must automatically prefill the Split Configuration inputs to target approximately **1,000 rows per file** using the following hardcoded logic:

* Set the default selected split mode to: **"Split into X Total Files"**.
* Calculate the optimal number of files ($X$) using the formula:

$$X = \max\left(1, \text{round}\left(\frac{\text{Total Rows}}{1000}\right)\right)$$


* Prefill the numeric input field with this calculated value of $X$.
* Automatically refresh the **Live Calculation Label** so the user immediately sees the theoretical resulting rows per file (which will naturally hover around $1000 \pm 100$).
* *The user can still manually overwrite this value if they choose.*

### Step 3: Hardcoded Schema Validation

Before any file splitting occurs, validate **every row** against a strict, hardcoded schema using struct tags.

```go
type Record struct {
    ID        int    `csv:"id" validate:"required,numeric"`
    Email     string `csv:"email" validate:"required,email"`
    Age       int    `csv:"age" validate:"gte=0,lte=120"`
    Status    string `csv:"status" validate:"oneof=active pending inactive"`
}

```

* **Error Handling:** If *any* row fails validation, abort the entire process. Collect all validation errors (e.g., `"Row 45: Email is invalid"`) and print them out clearly in the UI Log Terminal. Do not generate output files if validation fails.

### Step 4: Splitting & Exporting

If validation passes, split the data and write the output files to a **new subfolder** created in the same directory as the source file (named `split_output_YYYYMMDD_HHMMSS`).

* **File Naming:** Use the prefix from the UI input field followed by a hyphen and the incrementing index. If the prefix is `custom-name`, files must be named `custom-name-1.csv`, `custom-name-2.csv`, etc.
* **Mode A (Lines per file):** If the user inputs $1000$, every output file contains exactly $1000$ lines (plus the original header), with the last file containing the remainder.
* **Mode B (Total files):** If the user inputs $5$, evenly distribute the rows across exactly $5$ files. The **last file gets fewer rows** if the division is uneven.
* **Output format:** Always `.csv`. Delimiter is always **semicolon (`;`)**. Encoding is **UTF-8 with BOM** for Excel compatibility on Windows.

---

## 4. Agent Instructions for Code Generation

* **Reactive UI State:** Ensure that changing the split mode radio buttons or typing inside the numeric input triggers a re-calculation function that safely handles dividing by zero or empty inputs without crashing.
* **Concurrency:** Run the file loading, validation, and splitting operations inside a separate goroutine so that the Fyne UI remains smooth and responsive during processing.
* **No Install Packaging:** Ensure the code does not rely on local C libraries or external assets so it can be cross-compiled cleanly using standard Go cross-compilation tools (`GOOS=windows`, `GOOS=darwin`, `GOOS=linux`).

---

## 5. CI/CD & Release

* **Platform:** GitHub Actions.
* **Trigger:** Push on any tag matching `v*` (e.g. `v1.0.0`, `v2.3.4`).
* **Tool:** [goreleaser](https://goreleaser.com) — standard Go release automation.
* **Targets (cross-compile):**
  - Linux   → `GOOS=linux`   `GOARCH=amd64`
  - Mac     → `GOOS=darwin`  `GOARCH=amd64`
  - Windows → `GOOS=windows` `GOARCH=amd64`
* **Archive format:** `.zip` for all platforms.
* **Binary naming:** `couic-v<version>-<Platform>.zip`
  - Examples: `couic-v1.2.3-Linux.zip`, `couic-v1.2.3-Mac.zip`, `couic-v1.2.3-Windows.zip`
* **Artifacts:** Each archive contains the single executable. Checksums auto-generated.
* **Release:** goreleaser creates a GitHub Release and uploads all `.zip` archives + checksums.
* **Go version:** `1.22` (pinned). No C dependencies or external assets required.

---

## 6. Testing

### 6.1 Unit Tests (core logic — no UI dependency)

| Package | File | What to test |
|---|---|---|
| `core/validator_test.go` | Validate each field of `Record` (`id`, `email`, `age`, `status`). Abort on first invalid row. |
| `core/loader_test.go` | CSV/Excel loading, header detection (≥2 matches → header), delimiter auto-detect (`,` vs `;`), case-insensitive mapping, error on missing file. |
| `core/splitter_test.go` | Even distribution (total-files mode), rows-per-file mode, uneven last chunk, output dir naming. |
| `export/exporter_test.go` | Semicolon delimiter, UTF-8 BOM, proper CSV quoting, empty file edge-case. |

### 6.2 UI Tests

Use the Fyne `test` package (`fyne.io/fyne/v2/test`):
- Widget state transitions: header checkbox toggling, mode radio switch, file pick OK/cancel.
- Reactive label updates on input changes.
- Process button enabled/disabled state.

### 6.3 Test Fixtures

- Small valid CSV with `,` and `;` delimiters.
- Small valid Excel (`.xlsx`).
- CSV with headers / without headers.
- Invalid data (missing fields, bad email, non-numeric age).
- Empty file.

### 6.4 Run Command

```sh
go test ./... -v
```
