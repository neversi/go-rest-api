package service

import (
	"fmt"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
)

// TaskController funcs -> TaskService

// ITaskService interface is needed to be able to flexible connection of different ports
type ITaskService interface {
	Create(t *models.Task) 	error
	Read(t *models.Task) 	([]*models.Task, error) 
	Update(t *models.Task) 	error 
	Delete(t *models.Task) 	error
	FindByUserID(id uint) 	([]*models.Task, error)
	FindByID(id uint) 	(*models.Task, error)
}

// TaskService of service
type TaskService struct {
	TaskRepository *database.TaskRepository
}

// Create creates the function
func (service *TaskService) Create(t *models.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}

	err := service.TaskRepository.Create(t)
	if err != nil {
		return err
	}
	return nil
}

// Read ...
func (service *TaskService) Read(t *models.Task) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)
	tasks, err := service.TaskRepository.Read(t)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Update ...
func (service *TaskService) Update(t *models.Task) error {
	task, err := service.FindByUserID(t.ID);

	if err != nil {
		return err
	}
	
	if task == nil {
		return fmt.Errorf("There is no such task");
	}
	
	err = service.TaskRepository.Update(t)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (service *TaskService) Delete(t *models.Task) error {
	task, err := service.FindByID(t.ID);

	if err != nil {
		return err
	}

	if task == nil {
		return fmt.Errorf("There is no such task");
	}
	
	err = service.TaskRepository.Delete(t)

	if err != nil { 
		return err
	}

	return nil
}



// FindByUserID ...
func (service *TaskService) FindByUserID(id uint) ([]*models.Task, error) {
	task := new(models.Task)
	task.UserID = id
	tasks, err := service.TaskRepository.Read(task)

	if err != nil {
		return nil, err
	}
	
	return tasks, nil
}

// FindByID ...
func (service *TaskService) FindByID(id uint) (*models.Task, error) {
	task := new(models.Task)
	task.ID = id
	tasks, err := service.TaskRepository.Read(task)

	if err != nil {
		return nil, err
	}

	return tasks[0], nil
}








