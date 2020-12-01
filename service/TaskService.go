package service

import (
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
)

// TaskController funcs -> TaskService

// ITaskService interface is needed to be able to flexible connection of different ports
type ITaskService interface {
	Create(t *models.Task) 	error
	Read(t *models.Task) 	error 
	Update(t *models.Task) 	error 
	Delete(t *models.Task) 	error
	FindByID(id uint) 	(*models.Task)
}

// TaskService of service
type TaskService struct {
	TaskRepository *database.TaskRepository
}

// Create creates the function
func (ts *TaskService) Create(t *models.Task) error {

	err := ts.TaskRepository.Create(t)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (ts *TaskService) Delete(t *models.Task) error {
	return nil
}

// Update ...
func (ts *TaskService) Update(t *models.Task) error {
	return nil
}

// Read ...
func (ts *TaskService) Read(t *models.Task) error {
	return nil
}

// FindByID ...
func (ts *TaskService) FindByID(id uint) *models.Task {
	return nil
}







