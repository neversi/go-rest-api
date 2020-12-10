package models

import validation "github.com/go-ozzo/ozzo-validation/v4"

// Role ...
type Role struct {
	ID 	uint 	"json:\"id\""
	Role	string	"json:\"role\" gorm:\"default: user\""
	UserID	uint	"json:\"user_id\""
	User 	User	
}

// Validate record to the right place
func (r *Role) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Role, validation.Empty),
		validation.Field(&r.UserID, validation.NotNil, validation.Empty),
	)
}
