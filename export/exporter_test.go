package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"couic/core"
)

func TestWriteCSV_SemicolonDelimiter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.csv")

	recs := []core.Record{
		{PrenomUtilisateur: "Jean", Email: "jean@example.com", Categorie: "A"},
	}
	colOrder := []int{0, 3, 5} // PrenomUtilisateur, Email, Categorie

	if err := WriteCSV(path, recs, colOrder); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if !strings.HasPrefix(content, "\xef\xbb\xbf") {
		t.Fatal("missing UTF-8 BOM")
	}

	body := content[3:]
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (header + data), got %d", len(lines))
	}

	header := strings.Split(lines[0], ";")
	if len(header) != 3 {
		t.Fatalf("expected 3 columns, got %d: %v", len(header), header)
	}
	if header[0] != "Prénom de l'utilisateur" {
		t.Errorf("expected Prénom de l'utilisateur, got %s", header[0])
	}

	dataRow := strings.Split(lines[1], ";")
	if dataRow[0] != "Jean" {
		t.Errorf("expected Jean, got %s", dataRow[0])
	}
	if dataRow[1] != "jean@example.com" {
		t.Errorf("expected jean@example.com, got %s", dataRow[1])
	}
}

func TestWriteCSV_EmptyRecords(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.csv")

	if err := WriteCSV(path, nil, []int{0, 3}); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	body := string(data[3:])
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line (header only), got %d", len(lines))
	}
}

func TestWriteCSV_ColumnOrderSourceOrder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.csv")

	recs := []core.Record{
		{PrenomUtilisateur: "Jean", NomUtilisateur: "Dupont", Email: "jean@example.com"},
	}
	// Source had columns: Email, PrenomUtilisateur, NomUtilisateur
	colOrder := []int{3, 0, 1}

	if err := WriteCSV(path, recs, colOrder); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	body := string(data[3:])
	lines := strings.Split(strings.TrimSpace(body), "\n")

	header := strings.Split(lines[0], ";")
	if header[0] != "Email" {
		t.Errorf("expected Email first, got %s", header[0])
	}
	if header[1] != "Prénom de l'utilisateur" {
		t.Errorf("expected Prénom second, got %s", header[1])
	}

	dataRow := strings.Split(lines[1], ";")
	if dataRow[0] != "jean@example.com" {
		t.Errorf("expected jean@example.com, got %s", dataRow[0])
	}
}

func TestWriteAll_CreatesFiles(t *testing.T) {
	dir := t.TempDir()

	recs := []core.Record{
		{PrenomUtilisateur: "A"},
		{PrenomUtilisateur: "B"},
		{PrenomUtilisateur: "C"},
	}
	chunks := [][]core.Record{recs[0:1], recs[1:3]}

	paths, err := WriteAll(dir, "split", chunks, []int{0})
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 2 {
		t.Fatalf("expected 2 files, got %d", len(paths))
	}

	for i, p := range paths {
		if !filepath.IsAbs(p) {
			t.Errorf("expected absolute path, got %s", p)
		}
		expected := filepath.Join(dir, "split")
		if !strings.HasPrefix(p, expected) {
			t.Errorf("expected path to start with %s, got %s", expected, p)
		}
		if !strings.HasSuffix(p, ".csv") {
			t.Errorf("expected .csv extension, got %s", p)
		}
		_ = i
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries in output dir, got %d", len(entries))
	}
}
