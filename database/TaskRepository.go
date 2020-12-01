package database

import (
	"encoding/json"
	"fmt"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gorm.io/gorm"
)

// TaskRepository ...
type TaskRepository struct {
	db *DataBase
}

// NewTaskRepository creates an instance of TaskRepository
func NewTaskRepository(db *DataBase) *TaskRepository {
	return &TaskRepository{db: db}
}
// Create the task
func (tr *TaskRepository) Create(t *models.Task) (error) {
	currentDB := tr.db.Pdb

	if err := t.Validate(); err != nil {
		return err
	}

	currentDB.Create(t)
	return nil
}

// Read retrieves all tasks related to t
func (tr *TaskRepository) Read(t *models.Task) ([]*models.Task, error) {
	currentDB := tr.db.Pdb
	tasks := make([]*models.Task, 0)
	var result *gorm.DB
	if (t == nil) {
		result = currentDB.Table("tasks").Select("*").Find(&tasks)
	} else {
		result = currentDB.Model(&models.Task{}).Where("user_id = ?", t.UserID).Find(&tasks)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

// Delete the task
func (tr *TaskRepository) Delete(t *models.Task) error {
	currentDB := tr.db.Pdb

	result := currentDB.Where("id = ?", t.ID).Delete(&models.Task{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update updates the info about user
func (tr *TaskRepository) Update(t *models.Task) (*models.Task, error) {
	currentDB := tr.db.Pdb

	var oldObj = new(models.Task)

	result := currentDB.Table("tasks").Where("id = ?", t.ID).First(&oldObj)

	
	if result.Error != nil {
		return nil, result.Error
	}

	updateMap:= make(map[string][]byte)
	oldMap:= make(map[string][]byte)
	newMap:= make(map[string][]byte)

	oldBytes, err := json.Marshal(oldObj)
	_ = json.Unmarshal(oldBytes, &oldMap)
	newBytes, err := json.Marshal(t)
	_ = json.Unmarshal(newBytes, &updateMap)
	fmt.Println(updateMap, oldMap)

	for key, value := range updateMap {
		
		if len(value) > 0 {
			newMap[key] = updateMap[key]
		} else {
			newMap[key] = oldMap[key]
		}
	}

	newBytes, err = json.Marshal(newMap)
	_ = json.Unmarshal(newBytes, &oldObj)
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if err = tr.db.Pdb.Model(&models.Task{}).Where("id = ?", t.ID).Save(&oldObj).Error;
	err != nil {
		return nil, err
	}
	
	return oldObj, nil
}