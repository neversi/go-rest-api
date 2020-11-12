package objects

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Record of the author (zadacha)
type Record struct {
	Title    string "json:\"title\""
	Text     string "json:\"text\""
	Category string "json:\"category\""
	Author   string "json:\"author\""
	Status   string "json:\"status\""
	Deadline string "json:\"date\""
}

// Validation of record
func (currentRecord *Record) Validation() error {
	return validation.ValidateStruct(
		currentRecord,
		validation.Field(&currentRecord.Title, validation.Length(6, 0)),
		validation.Field(&currentRecord.Text, validation.Length(1, 0)),
	)
}

