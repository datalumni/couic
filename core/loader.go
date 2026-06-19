package core

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
)

type LoadResult struct {
	Records     []Record
	ColumnOrder []int
	HasHeader   bool
	TotalRows   int
}

func LoadFile(path string) (*LoadResult, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".csv":
		return loadCSV(path)
	case ".xlsx":
		return loadExcel(path)
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

func loadExcel(path string) (*LoadResult, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir le fichier Excel : %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("le classeur ne contient aucune feuille")
	}
	if len(sheets) > 1 {
		fmt.Printf("Attention : le classeur contient %d feuilles. Seule la première (%s) sera lue.\n", len(sheets), sheets[0])
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("erreur de lecture de la feuille %s : %w", sheets[0], err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("la feuille est vide")
	}

	rawRows := make([][]string, len(rows))
	for i, row := range rows {
		rawRows[i] = row
	}

	return rowsToRecords(rawRows)
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
	for _, col := range headerRow {
		lower := strings.ToLower(col)
		found := false
		for tag, idx := range tagToIdx {
			if strings.ToLower(tag) == lower {
				order = append(order, idx)
				found = true
				break
			}
		}
		if !found {
			order = append(order, -1)
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
		if fieldIdx < 0 || j >= len(row) {
			continue
		}
		rv.Field(fieldIdx).SetString(strings.TrimSpace(row[j]))
	}
	return rec, nil
}
