package core

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ValidationError struct {
	Row   int
	Field string
	Tag   string
}

func (e ValidationError) String() string {
	msg := frenchError(e.Tag)
	return fmt.Sprintf("Ligne %d : %s : %s", e.Row, e.Field, msg)
}

func frenchError(tag string) string {
	switch tag {
	case "required":
		return "champ obligatoire"
	case "min":
		return "valeur trop courte"
	case "email":
		return "email invalide"
	case "datetime":
		return "format de date invalide (JJ/MM/AAAA attendu)"
	default:
		return "valeur invalide"
	}
}

func ValidateRecords(records []Record) []ValidationError {
	nameToLabel := make(map[string]string)
	for _, fi := range Fields() {
		nameToLabel[fi.Name] = fi.CSVTag
	}

	for i, rec := range records {
		errs := validate.Struct(rec)
		if errs == nil {
			continue
		}

		var out []ValidationError
		for _, verr := range errs.(validator.ValidationErrors) {
			fieldName := nameToLabel[verr.StructField()]
			if fieldName == "" {
				fieldName = verr.Field()
			}
			out = append(out, ValidationError{
				Row:   i + 1,
				Field: fieldName,
				Tag:   verr.Tag(),
			})
		}
		if len(out) > 0 {
			return out
		}
	}
	return nil
}
