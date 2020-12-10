package service

import (
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
)

// IRoleService ...
type IRoleService interface {
	FindByUserID(userID uint) (string, error)
		SetUserRole(userID uint, role string) error 
}

// RoleService ...
type RoleService struct {
	RoleRepository *database.RoleRepository
}

// FindByUserID ...
func (service *RoleService) FindByUserID(userID uint) (string, error) {
	userRole := new(models.Role)
	userRole.UserID = userID
	userRole, err := service.RoleRepository.Read(userRole)

	if err != nil {
		return "", err
	}
	
	return userRole.Role, nil
}

// SetUserRole ... 
func (service *RoleService) SetUserRole(userID uint, role string) error {
	if role == "" {
		return nil
	}
	userRole := new(models.Role)
	userRole.UserID = userID
	userRole.Role = role
	
	err := service.RoleRepository.Update(userRole)
	if err != nil {
		return nil
	}
	
	return nil
}