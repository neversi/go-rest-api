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
	FindByID(id uint) 	(*models.Task, error)
}

// TaskService of service
type TaskService struct {
	TaskRepository *database.TaskRepository
}

// Create creates the function
func (ts *TaskService) Create(t *models.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}

	err := ts.TaskRepository.Create(t)
	if err != nil {
		return err
	}
	return nil
}

// Read ...
func (ts *TaskService) Read(t *models.Task) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)
	tasks, err := ts.TaskRepository.Read(t)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Update ...
func (ts *TaskService) Update(t *models.Task) error {
	task, err := ts.FindByID(t.ID);

	if err != nil {
		return err
	}
	
	if task == nil {
		return fmt.Errorf("There is no such task");
	}
	
	err = ts.TaskRepository.Update(t)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (ts *TaskService) Delete(t *models.Task) error {
	task, err := ts.FindByID(t.ID);

	if err != nil {
		return err
	}

	if task == nil {
		return fmt.Errorf("There is no such task");
	}
	
	err = ts.TaskRepository.Delete(t)

	if err != nil { 
		return err
	}

	return nil
}



// FindByID ...
func (ts *TaskService) FindByID(id uint) (*models.Task, error) {
	task := new(models.Task)
	task.ID = id
	tasks, err := ts.TaskRepository.Read(task)

	if err != nil {
		return nil, err
	}
	
	return tasks[0], nil
}







