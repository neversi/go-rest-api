package database

import (
	// "fmt"

	"encoding/json"

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
	users := make([]*models.User, 0)
	var result *gorm.DB
	if u == nil {
		result = currentDB.Table("users").Select("*").Find(&users)
	} else {
		result = currentDB.Where("login = ?", u.Login).First(&users)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	// sqlResult, err := result.Rows()
	// if err != nil {
	// 	return nil, err
	// }

	// for sqlResult.Next() {
	// 	user := new(models.User)
	// 	err := sqlResult.Scan(&user.ID, &user.Login, &user.Password, &user.FirstName, &user.SurName, &user.Email)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	users = append(users, user)
	// }
	// if err = sqlResult.Err(); err != nil {
	// 	return nil, err
	// }
	return users, nil
}

// Update updates the info about user
func (ur *UserRep) Update(u *models.User) (*models.User, error) {
	currentDB := ur.DB.Pdb

	var user = new(models.User)
	var result *gorm.DB
	if result = currentDB.Table("users").Select("*").Where("login = ?", u.Login).First(&models.User{}); result.Error != nil {
		return nil, result.Error
	}
	res := result.Row()
	res.Scan(&user.ID, &user.Login, &user.Password, &user.FirstName, &user.SurName, &user.Email)

	
	ur.DB.Lock()
	var updates, oldMap map[string]interface{}
	oldRec, _ := json.Marshal(user)
	_ = json.Unmarshal(oldRec, &oldMap)
	newRec, _ := json.Marshal(u)
	_ = json.Unmarshal(newRec, &updates)
	changedPassword := false;
	for key, value := range updates {
		if key == "password" && value != "" {
			changedPassword = true;
		}
		if value == "" {
			updates[key] = oldMap[key]
		}
	}
	newRec, _ = json.Marshal(updates)
	_ = json.Unmarshal(newRec, &u)

	if err := u.Validate(); err != nil {
		return nil, err
	}
	
	if changedPassword == true {
		encrypted, err := EncryptString(updates["password"].(string))
		if err != nil {
			return nil, err
		}
		u.Password = encrypted
	}

	ur.DB.Pdb.Table("users").Where("login = ? ", u.Login).Save(&u)

	ur.DB.Unlock()

	return u, nil
}

// Delete user by the login parameter
func (ur *UserRep) Delete(u *models.User) error {
	currentDB := ur.DB.Pdb

		ur.DB.Lock()
		currentDB.Where("login = ?", u.Login).Delete(&models.User{})
	ur.DB.Unlock()

	return nil
}

