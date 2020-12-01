package service

import (
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
)

// IUserService ...
type IUserService interface {
	Create(u *models.User) 		error
	Read(u *database.UserDTO) 	([]*database.UserDTO, error)
	Update(u *database.UserDTO)	error
	Delete(u *database.UserDTO)	error 
	FindByLogin(login string)	(*database.UserDTO, error)
	FindByID(id uint)		(*database.UserDTO, error)
}

// UserService ...
type UserService struct {
	UserRepository *database.UserRepository
}

// Create ...
func (service *UserService) Create(usr *models.User) error {
	user, err := service.FindByLogin(usr.Login)

	if err = usr.Validate(); err != nil {
		return err
	}

	if 
	
}

// Delete ...
func (service *UserService) Delete(t *database.UserDTO) error {
	return nil
}

// Update ...
func (service *UserService) Update(t *database.UserDTO) error {

	return nil
}

// Read ...
func (service *UserService) Read(t *database.UserDTO) ([]*database.UserDTO, error) {

	users, err := service.UserRepository.Read(t)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// FindByLogin ...
func (service *UserService) FindByLogin(login string) (*models.User, error) {
	user := new(models.User)
	user.Login = login
	users, err := service.UserRepository.Read(user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return user[0]
}

