package database

import (
	// "fmt"

	"encoding/json"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gorm.io/gorm"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

)

// UserRepository user's manager
type UserRepository struct {
	db *DataBase
}

// UserDTO is the struct to transfer objects between different layers
type UserDTO struct {
	ID		uint	"json:\"id\""
	Login 		string	"json:\"login\""
	FirstName	string	"json:\"first_name\""	
	SurName 	string	"json:\"sur_name\""
	Email 		string	"json:\"email\""
}

// Validate validates the user information before creating or upgrading it
func (u *UserDTO) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Login, validation.Length(6, 30), is.Alpha),
		validation.Field(&u.FirstName, validation.Length(2, 40), is.Alpha),
		validation.Field(&u.SurName, validation.Length(2, 40), is.Alpha),
		validation.Field(&u.Email, is.EmailFormat),
	)
}

// NewUserRepository creates an instance of UserRepository
func NewUserRepository(db *DataBase) *UserRepository {
	return &UserRepository{db: db}
}

// Create ...
func (ur *UserRepository) Create(u *models.User) error {

	currentDB := ur.db.Pdb
	encrypted, err := misc.EncryptString(u.Password)
	if err != nil {
		return err
	}

	u.Password = encrypted
	currentDB.Create(u)
	return nil
}

// Read retrieves the user(s) from table users
func (ur *UserRepository) Read(u *models.User) ([]*models.User, error) {
	currentDB := ur.db.Pdb
	users := make([]*models.User, 0)
	var result *gorm.DB
	if u == nil {
		result = currentDB.Model(&models.User{}).Select("*").Find(&users)
	} else {
		if len(u.Login) == 0 {
			result = currentDB.Model(&models.User{}).Where("id = ?", u.ID).First(&users)
		} else {
			result = currentDB.Model(&models.User{}).Where("login = ?", u.Login).First(&users)
		}
	}

	if result.Error != nil && len(users) != 0 {
		return nil, result.Error
	}

	if len(users) == 0 {
		return nil, nil
	}
	
	return users, nil
}

// Update updates the info about user
func (ur *UserRepository) Update(u *models.User) error {
	currentDB := ur.db.Pdb

	var user = new(models.User)
	var result *gorm.DB
	if result = currentDB.Table("users").Select("*").Where("id = ?", u.ID).First(&user); result.Error != nil {
		return result.Error
	}

	if err := u.Validate(); err != nil {
		return err
	}

	oldMap := make(map[string][]byte)
	changes := make(map[string][]byte)

	oldBytes, _ := json.Marshal(user)
	_ = json.Unmarshal(oldBytes, &oldMap)

	changedBytes, _ := json.Marshal(u)
	_ = json.Unmarshal(changedBytes, &changes)
	
	for key, value := range changes {
		if key == "login" || key == "id" {
			continue
		}
		if len(value) > 0 {
			oldMap[key] = value
		}
	}

	oldBytes, _ = json.Marshal(oldMap)
	_ = json.Unmarshal(oldBytes, &user)
	
	err := ur.db.Pdb.Model(&models.User{}).Where("login = ? ", u.Login).Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete user by the login parameter
func (ur *UserRepository) Delete(u *models.User) error {
	currentDB := ur.db.Pdb
	if (len(u.Login) != 0) {
		currentDB.Model(&models.User{}).Where("login = ?", u.Login).Delete(&models.User{})
	} else {
		currentDB.Model(&models.User{}).Where("id = ?", u.ID).Delete(&models.User{});
	}

	return nil
}

