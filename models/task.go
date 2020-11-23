package models

import (
	// "gorm.io/gorm"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Task of the author (zadacha)
type Task struct {
	ID 		uint   "json:\"id\" gorm:\"primaryKey\""
	UserRefer	uint   "json:\"user_id\""
	User		User	"gorm:\"foreignKey:UserRefer\""
	Title    	string "json:\"title\" gorm:\"not null;unique\""
	Text     	string "json:\"text\" gorm:\"not null\""
	Category 	string "json:\"category\" gorm:\"not null\""
	Author   	string "json:\"author\" gorm:\"not null\""
	Status   	string "json:\"status\" gorm:\"default:in progress\""
	Deadline 	string "json:\"date\""
}

// Validate record to the right place
func (currentRecord *Task) Validate() error {
	return validation.ValidateStruct(
		currentRecord,
		validation.Field(&currentRecord.Title, validation.Length(6, 0)),
		validation.Field(&currentRecord.Text, validation.Length(1, 0)),
	)
}
