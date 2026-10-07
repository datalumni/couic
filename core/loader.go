package core

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

var shortYearDateRe = regexp.MustCompile(`^(\d{1,2})([/-])(\d{1,2})[/-](\d{2})$`)

// expandShortYear turns a "01-06-26" style date into "01-06-2026".
// ponytail: 00-29 -> 20xx, 30-99 -> 19xx pivot; wrong for dates >30y old with a 2-digit year, adjust pivot if that surfaces.
func expandShortYear(s string) string {
	m := shortYearDateRe.FindStringSubmatch(s)
	if m == nil {
		return s
	}
	yy, _ := strconv.Atoi(m[4])
	century := "19"
	if yy <= 29 {
		century = "20"
	}
	return m[1] + m[2] + m[3] + m[2] + century + m[4]
}

type LoadResult struct {
	Records     []Record
	ColumnOrder []int
	HasHeader   bool
	TotalRows   int
	Sheets      []string // Excel only: all sheet names
	Sheet       string   // Excel only: sheet actually read
}

// LoadFile loads a CSV or Excel file. For Excel, sheet picks the sheet to read ("" = first).
func LoadFile(path, sheet string) (*LoadResult, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".csv":
		return loadCSV(path)
	case ".xlsx":
		return loadExcel(path, sheet)
	default:
		return nil, fmt.Errorf("format non supporté : %s", ext)
	}
}

func loadCSV(path string) (*LoadResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir le fichier : %w", err)
	}
	defer f.Close()

	sampled, err := sampleLines(f, 5)
	if err != nil {
		return nil, err
	}
	if len(sampled) == 0 {
		return nil, fmt.Errorf("fichier vide")
	}

	delim := detectDelimiter(sampled)

	f.Seek(0, 0)
	r := csv.NewReader(f)
	r.Comma = delim
	r.LazyQuotes = true
	r.FieldsPerRecord = -1

	rawRows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("erreur de lecture CSV : %w", err)
	}
	if len(rawRows) == 0 {
		return nil, fmt.Errorf("fichier vide")
	}

	return rowsToRecords(rawRows)
}

func loadExcel(path, sheet string) (*LoadResult, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir le fichier Excel : %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("le classeur ne contient aucune feuille")
	}
	if sheet == "" {
		sheet = sheets[0]
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("erreur de lecture de la feuille %s : %w", sheet, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("la feuille %s est vide", sheet)
	}

	rawRows := make([][]string, len(rows))
	for i, row := range rows {
		rawRows[i] = row
	}

	res, err := rowsToRecords(rawRows)
	if err != nil {
		return nil, err
	}
	res.Sheets = sheets
	res.Sheet = sheet
	return res, nil
}

func sampleLines(r io.Reader, n int) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() && len(lines) < n {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func detectDelimiter(samples []string) rune {
	commaCounts := make([]int, len(samples))
	semicolonCounts := make([]int, len(samples))
	for i, line := range samples {
		for _, ch := range line {
			if ch == ',' {
				commaCounts[i]++
			} else if ch == ';' {
				semicolonCounts[i]++
			}
		}
	}

	commaConsistent := allEqual(commaCounts)
	semicolonConsistent := allEqual(semicolonCounts)

	if commaConsistent && !semicolonConsistent {
		return ','
	}
	if semicolonConsistent && !commaConsistent {
		return ';'
	}

	sumComma := sum(commaCounts)
	sumSemi := sum(semicolonCounts)
	if sumSemi > sumComma {
		return ';'
	}
	return ','
}

func allEqual(vals []int) bool {
	if len(vals) <= 1 {
		return true
	}
	for i := 1; i < len(vals); i++ {
		if vals[i] != vals[0] {
			return false
		}
	}
	return true
}

func sum(vals []int) int {
	s := 0
	for _, v := range vals {
		s += v
	}
	return s
}

func rowsToRecords(rows [][]string) (*LoadResult, error) {
	tagToIdx := CSVTagToFieldIndex()

	hasHeader, headerRow := detectHeader(rows[0])
	var colOrder []int

	if hasHeader {
		colOrder = buildColumnOrder(headerRow, tagToIdx)
	} else {
		colOrder = defaultColumnOrder()
	}

	dataStart := 0
	if hasHeader {
		dataStart = 1
	}

	records := make([]Record, 0, len(rows)-dataStart)
	for i := dataStart; i < len(rows); i++ {
		rec, err := rowToRecord(rows[i], colOrder)
		if err != nil {
			return nil, fmt.Errorf("ligne %d : %w", i+1, err)
		}
		records = append(records, *rec)
	}

	return &LoadResult{
		Records:     records,
		ColumnOrder: colOrder,
		HasHeader:   hasHeader,
		TotalRows:   len(records),
	}, nil
}

func detectHeader(firstRow []string) (bool, []string) {
	normalized := make([]string, len(firstRow))
	for i, v := range firstRow {
		normalized[i] = strings.TrimSpace(v)
	}

	tags := CSVTags()
	matches := 0
	for _, col := range normalized {
		lower := strings.ToLower(col)
		for _, tag := range tags {
			if strings.ToLower(tag) == lower {
				matches++
				break
			}
		}
	}
	return matches >= 2, normalized
}

func buildColumnOrder(headerRow []string, tagToIdx map[string]int) []int {
	order := make([]int, 0, len(headerRow))
	for colIdx, col := range headerRow {
		lower := strings.ToLower(col)
		matched := false
		for tag, idx := range tagToIdx {
			if strings.ToLower(tag) == lower {
				order = append(order, idx)
				matched = true
				break
			}
		}
		if !matched {
			order = append(order, colIdx)
		}
	}
	return order
}

func defaultColumnOrder() []int {
	fields := Fields()
	order := make([]int, len(fields))
	for i := range fields {
		order[i] = i
	}
	return order
}

func rowToRecord(row []string, colOrder []int) (*Record, error) {
	rec := &Record{}
	rv := reflect.ValueOf(rec).Elem()
	for j, fieldIdx := range colOrder {
		if fieldIdx < 0 || fieldIdx >= rv.NumField() || j >= len(row) {
			continue
		}
		val := strings.TrimSpace(row[j])
		if IsDateField(fieldIdx) {
			val = expandShortYear(val)
		}
		rv.Field(fieldIdx).SetString(val)
	}
	return rec, nil
}
