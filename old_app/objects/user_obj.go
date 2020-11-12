package objects

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

)
// User is the struct of user which will use the app itself
type User struct {
	ID        string "json:\"id\""
	Login     string "json:\"login\""
	Password  string "json:\"password\""
	FirstName string "json:\"firstName\""
	SurName   string "json:\"surName\""
	Email     string "json:\"email\""
}

// Validation func validate the correctness of the data
func (currentUser *User) Validation() error {
	return validation.ValidateStruct(
		currentUser,
		validation.Field(&currentUser.Login, validation.Length(6, 30), is.Alpha),
		validation.Field(&currentUser.Password, validation.Length(8, 36), is.Alpha),
		validation.Field(&currentUser.FirstName, validation.Length(2, 40), is.Alpha),
		validation.Field(&currentUser.SurName, validation.Length(2, 40), is.Alpha),
		validation.Field(&currentUser.Email, is.EmailFormat),
	)
}