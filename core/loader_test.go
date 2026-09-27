package core

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempCSV(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadCSV_WithHeaderSemicolon(t *testing.T) {
	csv := "Prénom de l'utilisateur;Nom de l'utilisateur;Date de naissance;Email;Référence externe;Catégorie;Date de fin;Diplôme;Site;N° RNCP\n" +
		"Jean;Dupont;15/03/1990;jean@example.com;REF123;A;31/12/2025;Master;Paris;RNCP12345\n" +
		"Marie;Curie;;marie@example.com;;B;;;Lyon;\n"
	path := writeTempCSV(t, "test.csv", csv)

	res, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasHeader {
		t.Fatal("expected header detection")
	}
	if res.TotalRows != 2 {
		t.Fatalf("expected 2 rows, got %d", res.TotalRows)
	}
	if len(res.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(res.Records))
	}
	if res.Records[0].PrenomUtilisateur != "Jean" {
		t.Errorf("expected Jean, got %s", res.Records[0].PrenomUtilisateur)
	}
	if res.Records[0].Email != "jean@example.com" {
		t.Errorf("expected jean@example.com, got %s", res.Records[0].Email)
	}
	if res.Records[1].PrenomUtilisateur != "Marie" {
		t.Errorf("expected Marie, got %s", res.Records[1].PrenomUtilisateur)
	}
	if res.Records[1].Categorie != "B" {
		t.Errorf("expected B, got %s", res.Records[1].Categorie)
	}
}

func TestLoadCSV_WithHeaderComma(t *testing.T) {
	csv := "Prénom de l'utilisateur,Nom de l'utilisateur,Email,Date de naissance,Catégorie,Référence externe,Date de fin,Diplôme,Site,N° RNCP\n" +
		"Jean,Dupont,jean@example.com,15/03/1990,A,REF123,31/12/2025,Master,Paris,RNCP12345\n"
	path := writeTempCSV(t, "test.csv", csv)

	res, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasHeader {
		t.Fatal("expected header detection")
	}
	if res.TotalRows != 1 {
		t.Fatalf("expected 1 row, got %d", res.TotalRows)
	}
}

func TestLoadCSV_NoHeader(t *testing.T) {
	csv := "Jean;Dupont;15/03/1990;jean@example.com;REF123;A;31/12/2025;Master;Paris;RNCP12345\n" +
		"Marie;Curie;;marie@example.com;;B;;;Lyon;\n"
	path := writeTempCSV(t, "test.csv", csv)

	res, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if res.HasHeader {
		t.Fatal("expected no header detection for data-like first row")
	}
	if res.TotalRows != 2 {
		t.Fatalf("expected 2 rows, got %d", res.TotalRows)
	}
}

func TestLoadCSV_EmptyFile(t *testing.T) {
	path := writeTempCSV(t, "empty.csv", "")
	_, err := LoadFile(path)
	if err == nil {
		t.Fatal("expected error for empty file")
	}
}

func TestLoadCSV_HeaderOnly(t *testing.T) {
	csv := "Prénom de l'utilisateur;Nom de l'utilisateur;Date de naissance;Email;Référence externe;Catégorie;Date de fin;Diplôme;Site;N° RNCP\n"
	path := writeTempCSV(t, "header.csv", csv)

	res, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasHeader {
		t.Fatal("expected header detection")
	}
	if res.TotalRows != 0 {
		t.Fatalf("expected 0 rows, got %d", res.TotalRows)
	}
}

func TestLoadCSV_CaseInsensitiveHeader(t *testing.T) {
	csv := "PRÉNOM DE L'UTILISATEUR;NOM DE L'UTILISATEUR;EMAIL;CATÉGORIE\n" +
		"Jean;Dupont;jean@example.com;A\n"
	path := writeTempCSV(t, "test.csv", csv)

	res, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasHeader {
		t.Fatal("expected header detection with case-insensitive matching")
	}
	if res.Records[0].PrenomUtilisateur != "Jean" {
		t.Errorf("expected Jean, got %s", res.Records[0].PrenomUtilisateur)
	}
	if res.Records[0].Email != "jean@example.com" {
		t.Errorf("expected jean@example.com, got %s", res.Records[0].Email)
	}
}

func TestLoadCSV_DelimiterAutoDetectSemicolon(t *testing.T) {
	csv := "Prénom;Nom;Email;Catégorie\nJean;Dupont;jean@example.com;A\n"
	path := writeTempCSV(t, "test.csv", csv)

	res, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasHeader {
		t.Fatal("expected header detection")
	}
}

func TestLoadCSV_UnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("a,b,c\n1,2,3\n"), 0644)
	_, err := LoadFile(path)
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}

func TestLoadCSV_ColumnOrderPreservesSourceOrder(t *testing.T) {
	csv := "Email;Prénom de l'utilisateur;Nom de l'utilisateur;Catégorie\n" +
		"jean@example.com;Jean;Dupont;A\n"
	path := writeTempCSV(t, "test.csv", csv)

	res, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// ColumnOrder should be [fieldIndex_of_Email, fieldIndex_of_Prenom, fieldIndex_of_Nom, fieldIndex_of_Categorie]
	if len(res.ColumnOrder) != 4 {
		t.Fatalf("expected 4 column order entries, got %d", len(res.ColumnOrder))
	}
	// Verify Email comes first in column order
	emailIdx := CSVTagToFieldIndex()["Email"]
	if res.ColumnOrder[0] != emailIdx {
		t.Errorf("expected Email field index %d first, got %d", emailIdx, res.ColumnOrder[0])
	}
}

func TestLoadCSV_SingleColumnDoesNotTriggerHeader(t *testing.T) {
	csv := "Juste un nom\nJean\nMarie\n"
	path := writeTempCSV(t, "test.csv", csv)

	res, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if res.HasHeader {
		t.Fatal("expected no header for single column that doesn't match")
	}
	if res.TotalRows != 3 {
		t.Fatalf("expected 3 data rows, got %d", res.TotalRows)
	}
}

func TestExpandShortYear(t *testing.T) {
	cases := map[string]string{
		"01-06-26":   "01-06-2026",
		"01/06/26":   "01/06/2026",
		"15/03/90":   "15/03/1990",
		"15/03/1990": "15/03/1990",
		"":           "",
		"not-a-date": "not-a-date",
	}
	for in, want := range cases {
		if got := expandShortYear(in); got != want {
			t.Errorf("expandShortYear(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLoadCSV_NormalizesShortYearDates(t *testing.T) {
	csv := "Prénom de l'utilisateur,Nom de l'utilisateur,Date de naissance,Email,Référence externe,Catégorie,Date de fin,Diplôme,Site,N° RNCP\n" +
		"Jean,Dupont,01-06-26,jean.dupont@example.com,REF001,A,31/12/25,Master,Paris,RNCP-001\n"
	path := writeTempCSV(t, "short_year.csv", csv)

	result, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec := result.Records[0]
	if rec.DateNaissance != "01-06-2026" {
		t.Errorf("DateNaissance = %q, want %q", rec.DateNaissance, "01-06-2026")
	}
	if rec.DateFin != "31/12/2025" {
		t.Errorf("DateFin = %q, want %q", rec.DateFin, "31/12/2025")
	}
}
