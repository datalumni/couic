package ui

import (
	"testing"

	"couic/core"
)

func TestFormatValidation_GroupsByRow(t *testing.T) {
	errs := []core.ValidationError{
		{Row: 2, Field: "Prénom", Tag: "required"},
		{Row: 2, Field: "Email", Tag: "email"},
		{Row: 5, Field: "Email", Tag: "email"},
	}
	want := "Validation — 2 ligne(s) en erreur sur 10 :\n" +
		"  Ligne 2 : Prénom (champ obligatoire) · Email (email invalide)\n" +
		"  Ligne 5 : Email (email invalide)"
	if got := formatValidation(errs, 10); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
