package database

import (
	// "fmt"

	"fmt"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gorm.io/gorm"
)

// UserRep user's manager
type UserRep struct {
	DB *DataBase
}

// NewUserRep creates an instance of UserRep
func newUserRep(DB *DataBase) *UserRep {
	return &UserRep{DB: DB}
}

// Create ...
func (ur *UserRep) Create(u *models.User) (*models.User, error) {

	currentDB := ur.DB.Pdb
	encrypted, err := EncryptString(u.Password)
	if err != nil {
		return nil, err
	}

	if err = u.Validate(); err != nil {
		return nil, err
	}
	
	u.Password = encrypted
	ur.DB.Lock()
	currentDB.Create(u)
	ur.DB.Unlock()
	return u, nil
}

// Read retrieves the user(s) from table users
func (ur *UserRep) Read(u *models.User) ([]*models.User, error) {
	currentDB := ur.DB.Pdb
	var result *gorm.DB
	if u == nil {
		result = currentDB.Table("users").Select("*")
	} else {
		result = currentDB.Where("login = ?", u.Login).Find(&models.User{})
	}
	sqlResult, err := result.Rows()
	if err != nil {
		return nil, err
	}
	users := make([]*models.User, 0)

	for sqlResult.Next() {
		user := new(models.User)
		err := sqlResult.Scan(&user.ID, &user.Login, &user.Password, &user.FirstName, &user.SurName, &user.Email)
		if err != nil {
			return nil, err
		}
		user.Password = "NULL"
		users = append(users, user)
	}
	if err = sqlResult.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// Update updates the info about user
func (ur *UserRep) Update(u *models.User) (*models.User, error) {
	currentDB := ur.DB.Pdb

	encrypted, err := EncryptString(u.Password)
	if err != nil {
		return nil, err
	}
	var user = new(models.User)
	var result *gorm.DB
	if result = currentDB.Table("users").Select("*").Where("login = ?", u.Login); result.Error != nil {
		return nil, result.Error
	}
	res := result.Row()
	res.Scan(&user.Login)
	fmt.Print(u.Email, " |", user.Login, "| ", u.Login)
	if err = u.Validate(); err != nil {
		return nil, err	
	}

	u.Password = encrypted

	var substitute *models.User
	
	ur.DB.Lock()
	currentDB.Table("users").Where("login = ?", u.Login).Update("email", u.Email)
	currentDB.Model(&models.User{}).Where("login = ?", u.Login).Update("first_name", u.FirstName)
	currentDB.Model(&models.User{}).Where("login = ?", u.Login).Update("sur_name", u.SurName)
	currentDB.Model(&models.User{}).Where("login = ?", u.Login).Update("password", u.Password)
	ur.DB.Unlock()

	currentDB.Where(u.ID).First(&models.User{})
	return substitute, nil
}

// Delete user by the login parameter
func (ur *UserRep) Delete(u *models.User) error {
	currentDB := ur.DB.Pdb

	ur.DB.Lock()
	currentDB.Where("login = ?", u.Login).Delete(&models.User{})
	ur.DB.Unlock()

	return nil
}

