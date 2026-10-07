package core

import (
	"testing"
)

func TestValidateRecords_ValidRow(t *testing.T) {
	rec := Record{
		PrenomUtilisateur: "Jean",
		NomUtilisateur:    "Dupont",
		DateNaissance:     "15/03/1990",
		Email:             "jean@example.com",
		ReferenceExterne:  "REF123",
		Categorie:         "A",
		DateFin:           "31/12/2025",
		Diplome:           "Master",
		Site:              "Paris",
		NumeroRNCP:        "RNCP12345",
	}
	errs := ValidateRecords([]Record{rec})
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateRecords_ValidRowWithOptionalsEmpty(t *testing.T) {
	rec := Record{
		PrenomUtilisateur: "Jean",
		NomUtilisateur:    "Dupont",
		Email:             "jean@example.com",
		Categorie:         "A",
	}
	errs := ValidateRecords([]Record{rec})
	if len(errs) != 0 {
		t.Fatalf("expected no errors for empty optional fields, got %v", errs)
	}
}

func TestValidateRecords_MissingRequired(t *testing.T) {
	rec := Record{
		PrenomUtilisateur: "",
		NomUtilisateur:    "Dupont",
		Email:             "jean@example.com",
		Categorie:         "A",
	}
	errs := ValidateRecords([]Record{rec})
	if len(errs) == 0 {
		t.Fatal("expected error for missing PrenomUtilisateur")
	}
	if errs[0].Row != 1 {
		t.Errorf("expected row 1, got %d", errs[0].Row)
	}
}

func TestValidateRecords_BadEmail(t *testing.T) {
	rec := Record{
		PrenomUtilisateur: "Jean",
		NomUtilisateur:    "Dupont",
		Email:             "not-an-email",
		Categorie:         "A",
	}
	errs := ValidateRecords([]Record{rec})
	if len(errs) == 0 {
		t.Fatal("expected error for bad email")
	}
}

func TestValidateRecords_BadDate(t *testing.T) {
	rec := Record{
		PrenomUtilisateur: "Jean",
		NomUtilisateur:    "Dupont",
		DateNaissance:     "2025-03-15",
		Email:             "jean@example.com",
		Categorie:         "A",
	}
	errs := ValidateRecords([]Record{rec})
	if len(errs) == 0 {
		t.Fatal("expected error for bad date format")
	}
}

func TestValidateRecords_ReportsAllInvalidRows(t *testing.T) {
	recs := []Record{
		{
			PrenomUtilisateur: "Jean",
			NomUtilisateur:    "Dupont",
			Email:             "jean@example.com",
			Categorie:         "A",
		},
		{
			PrenomUtilisateur: "",
			NomUtilisateur:    "",
			Email:             "bad",
			Categorie:         "",
		},
		{
			PrenomUtilisateur: "Marie",
			NomUtilisateur:    "Curie",
			Email:             "marie@example.com",
			Categorie:         "B",
		},
		{
			PrenomUtilisateur: "Paul",
			NomUtilisateur:    "Martin",
			Email:             "bad",
			Categorie:         "C",
		},
	}
	errs := ValidateRecords(recs)
	rows := map[int]bool{}
	for _, e := range errs {
		rows[e.Row] = true
	}
	if len(rows) != 2 || !rows[2] || !rows[4] {
		t.Fatalf("expected errors on rows 2 and 4, got %v", errs)
	}
}

func TestValidateRecords_MultipleFieldsOnOneRow(t *testing.T) {
	rec := Record{
		PrenomUtilisateur: "",
		NomUtilisateur:    "",
		Email:             "bad",
		Categorie:         "",
	}
	errs := ValidateRecords([]Record{rec})
	if len(errs) == 0 {
		t.Fatal("expected errors")
	}
	// row 2 (bad email + empty Categorie)
	if errs[0].Row != 1 {
		t.Errorf("expected row 1, got %d", errs[0].Row)
	}
}
