package models

import (
	// "gorm.io/gorm"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// User is the struct of user which will use the app itself
type User struct {
	ID	  uint	 "json:\"id\" gorm:\"primaryKey;type:uint\""
	Login     string "json:\"login\" gorm:\"not null;unique; <-:create\""
	Password  string "json:\"password\" gorm:\"not null\""
	FirstName string "json:\"first_name\"" 
	SurName   string "json:\"sur_name\""
	Email     string "json:\"email\" gorm:\"not null;unique; <-:create\""
}


// Validate validates the user information before creating or upgrading it
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Login, validation.Length(6, 30), is.Alphanumeric),
		validation.Field(&u.Password, validation.Length(8, 36), is.Alphanumeric),
		validation.Field(&u.FirstName, validation.Length(2, 40), is.Alpha),
		validation.Field(&u.SurName, validation.Length(2, 40), is.Alpha),
		validation.Field(&u.Email, is.EmailFormat),
	)
}