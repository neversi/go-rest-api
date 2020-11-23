package database

import "gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"

// TaskRep ...
type TaskRep struct {
	db *DataBase
}

// NewTaskRep creates an instance of TaskRep
func newTaskRep(db *DataBase) *TaskRep {
	return &TaskRep{db: db}
}
// Create the task
func (tr *TaskRep) Create(t *models.Task) (*models.Task, error) {
	currentDB := tr.db.Pdb

	if err := t.Validate(); err != nil {
		return nil, err
	}

	tr.db.Lock()
	currentDB.Create(t)
	tr.db.Unlock()
	return t, nil
}

// Delete the task
func (tr *TaskRep) Delete(id uint) error {
	currentDB := tr.db.Pdb

	tr.db.Lock()
	currentDB.Where("id = ?", id).Delete(&models.Task{})
	tr.db.Unlock()

	return nil
}

// Update updates the info about user
func (tr *TaskRep) Update(t *models.Task) (*models.Task, error) {
	currentDB := tr.db.Pdb

	if err := t.Validate(); err != nil {
		return nil, err	
	}

	var substitute *models.Task
	
	err := currentDB.Where("id = ?", t.ID).First(substitute).Error
	if err != nil {
		return nil, err 
	}

	tr.db.Lock()
	substitute = t
	tr.db.Unlock()
	
	return substitute, nil
}