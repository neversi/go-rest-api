package service

import (
	"fmt"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
)

// IUserService ...
type IUserService interface {
	Create(u *models.User) 	error
	Read(u *database.UserDTO) 			([]*database.UserDTO, error)
	Update(u *database.UserDTO)			error
	Delete(u *database.UserDTO)			error 
	FindByLogin(login string)			(*models.User, error)
	FindByID(id uint)				(*database.UserDTO, error)
	CheckUser(login string, password string) 	error
}

// UserService ...
type UserService struct {
	UserRepository *database.UserRepository
}

// Create ...
func (service *UserService) Create(usr *models.User) error {
	
	user, err := service.FindByLogin(usr.Login);

	if err != nil {
		return err
	}

	if user != nil {
		return fmt.Errorf("Such user exists")
	}

	service.UserRepository.Create(usr)
	return nil
}

// Delete ...
func (service *UserService) Delete(t *database.UserDTO) error {
	user := new(models.User)
	user.Login = t.Login
	user.ID = t.ID
	tempUser := new(database.UserDTO)
	if len(t.Login) != 0 {
		temp, _ := service.FindByLogin(user.Login);
		tempUser = convertUserToDTO(temp);
	} else {
		tempUser, _ = service.FindByID(t.ID);
	}
	if tempUser == nil {
		return fmt.Errorf("User Does not exist");
	}
	service.UserRepository.Delete(user);
	return nil
}

// Update ...
func (service *UserService) Update(t *database.UserDTO) error {
	
	user := new(models.User)
	user.Login = t.Login;
	user.ID = t.ID;
	user.FirstName = t.FirstName;
	user.SurName = t.SurName;
	user.Email = t.Email;
	err := service.UserRepository.Update(user);
	if err != nil {
		return err;
	}
	return nil
}

// Read ...
func (service *UserService) Read(t *database.UserDTO) ([]*database.UserDTO, error) {
	var err error
	user := new(models.User)
	if t == nil {
		user = nil
	} else {
		user.Login = t.Login;
		user.ID = t.ID;
	}

	users := make([]*models.User, 0)
	
	users, err = service.UserRepository.Read(user);

	if err != nil {
		return nil, err
	}
	
	userDTOs := make([]*database.UserDTO, 0);

	for i := 0; i < len(users); i++ {
		userDTOs = append(userDTOs, convertUserToDTO(users[i]));
	}
	
	return userDTOs, nil
}

// FindByLogin ...
func (service *UserService) FindByLogin(login string) (*models.User, error) {
	user := new(models.User)
	user.Login = login
	users, err := service.UserRepository.Read(user)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

// CheckUser ...
func (service *UserService) CheckUser(login string, password string) error {
	user, err := service.FindByLogin(login)
	if err != nil {
		return err;
	}
	if user != nil {
		return fmt.Errorf("There is no such user");
	}
	exist, err := misc.CompareHash(user.Password, password);
	if err != nil {
		return err;
	}
	if exist == false {
		return fmt.Errorf("Incorrect password");
	}
	return nil;
}

// FindByID ... 
func (service *UserService) FindByID(id uint) (*database.UserDTO, error) {
	user := new(models.User)

	users, err := service.UserRepository.Read(user);

	if err != nil {
		return nil, err
	}
	
	if len(users) == 0 {
		return nil, nil
	}
	
	userDTO := convertUserToDTO(users[0]);

	return userDTO, nil
}

func convertUserToDTO(u *models.User) *database.UserDTO {
	if (u == nil) {
		return nil
	} 
	tempDTO := new(database.UserDTO);
	tempDTO.Email = u.Email;
	tempDTO.FirstName = u.FirstName;
	tempDTO.ID = u.ID;
	tempDTO.Login = u.Login;
	tempDTO.SurName = u.SurName;

	return tempDTO
}