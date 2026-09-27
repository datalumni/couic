package core

import (
	"reflect"
	"strings"
)

type Record struct {
	PrenomUtilisateur string `csv:"Prénom de l'utilisateur" validate:"required,min=1"`
	NomUtilisateur    string `csv:"Nom de l'utilisateur" validate:"required,min=1"`
	DateNaissance     string `csv:"Date de naissance" validate:"omitempty,datetime=02/01/2006"`
	Email             string `csv:"Email" validate:"required,email"`
	ReferenceExterne  string `csv:"Référence externe" validate:"omitempty,min=1"`
	Categorie         string `csv:"Catégorie" validate:"required,min=1"`
	DateFin           string `csv:"Date de fin" validate:"omitempty,datetime=02/01/2006"`
	Diplome           string `csv:"Diplôme" validate:"omitempty,min=1"`
	Site              string `csv:"Site" validate:"omitempty,min=1"`
	NumeroRNCP        string `csv:"N° RNCP" validate:"omitempty,min=1"`
}

type FieldInfo struct {
	Index  int
	CSVTag string
	Name   string
	IsDate bool
}

var recordFields []FieldInfo
var dateFieldIndices map[int]bool

func init() {
	t := reflect.TypeOf(Record{})
	recordFields = make([]FieldInfo, 0, t.NumField())
	dateFieldIndices = make(map[int]bool)
	for i := range t.NumField() {
		f := t.Field(i)
		isDate := strings.Contains(f.Tag.Get("validate"), "datetime=")
		recordFields = append(recordFields, FieldInfo{
			Index:  i,
			CSVTag: f.Tag.Get("csv"),
			Name:   f.Name,
			IsDate: isDate,
		})
		if isDate {
			dateFieldIndices[i] = true
		}
	}
}

func IsDateField(fieldIdx int) bool {
	return dateFieldIndices[fieldIdx]
}

func Fields() []FieldInfo {
	return recordFields
}

func CSVTags() []string {
	tags := make([]string, len(recordFields))
	for i, fi := range recordFields {
		tags[i] = fi.CSVTag
	}
	return tags
}

func CSVTagToFieldIndex() map[string]int {
	m := make(map[string]int, len(recordFields))
	for _, fi := range recordFields {
		m[fi.CSVTag] = fi.Index
	}
	return m
}
