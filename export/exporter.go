package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"couic/core"
)

func WriteCSV(path string, records []core.Record, columnOrder []int) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier : %w", err)
	}
	defer f.Close()

	_, err = f.WriteString("\xef\xbb\xbf")
	if err != nil {
		return fmt.Errorf("erreur d'écriture BOM : %w", err)
	}

	w := csv.NewWriter(f)
	w.Comma = ';'

	fields := core.Fields()

	if len(records) > 0 || len(columnOrder) > 0 {
		header := make([]string, len(columnOrder))
		for i, fi := range columnOrder {
			if fi >= 0 && fi < len(fields) {
				header[i] = fields[fi].CSVTag
			}
		}
		if err := w.Write(header); err != nil {
			return fmt.Errorf("erreur d'écriture d'en-tête : %w", err)
		}
	}

	for _, rec := range records {
		row := make([]string, len(columnOrder))
		rv := reflect.ValueOf(rec)
		for i, fi := range columnOrder {
			if fi >= 0 && fi < rv.NumField() {
				row[i] = rv.Field(fi).String()
			}
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("erreur d'écriture de ligne : %w", err)
		}
	}

	w.Flush()
	return w.Error()
}

func WriteAll(outputDir, prefix string, chunks [][]core.Record, columnOrder []int) ([]string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("impossible de créer le dossier %s : %w", outputDir, err)
	}

	var paths []string
	for i, chunk := range chunks {
		filename := fmt.Sprintf("%s-%d.csv", prefix, i+1)
		path := filepath.Join(outputDir, filename)
		if err := WriteCSV(path, chunk, columnOrder); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}
