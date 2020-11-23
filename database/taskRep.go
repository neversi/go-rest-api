package database

import (
	"encoding/json"
	"fmt"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gorm.io/gorm"
)

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

// Read retrieves all tasks related to t
func (tr *TaskRep) Read(t *models.Task) ([]*models.Task, error) {
	currentDB := tr.db.Pdb
	tasks := make([]*models.Task, 0)
	var result *gorm.DB
	if (t == nil) {
		result = currentDB.Table("tasks").Select("*").Find(&tasks)
	} else {
		result = currentDB.Model(&models.Task{}).Where("user_refer = ?", t.UserRefer).Find(&tasks)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

// Delete the task
func (tr *TaskRep) Delete(t *models.Task) error {
	currentDB := tr.db.Pdb

	tr.db.Lock()
	currentDB.Where("id = ?", t.ID).Delete(&models.Task{})
	tr.db.Unlock()

	return nil
}

// Update updates the info about user
func (tr *TaskRep) Update(t *models.Task) (*models.Task, error) {
	currentDB := tr.db.Pdb

	var oldObj = new(models.Task)

	result := currentDB.Table("tasks").Where("id = ?", t.ID).First(&models.Task{})

	
	if result.Error != nil {
		return nil, result.Error
	}

	var updateMap, oldMap map[string]interface{}
	oldBytes, err := json.Marshal(oldObj)
	_ = json.Unmarshal(oldBytes, &oldMap)
	newBytes, err := json.Marshal(t)
	_ = json.Unmarshal(newBytes, &updateMap)

	for key, value := range updateMap {
		if value.(string) != "" {
			fmt.Println("Here")
			oldMap[key] = updateMap[key]
		}
	}

	newBytes, err = json.Marshal(oldMap)
	_ = json.Unmarshal(newBytes, &oldObj)
	fmt.Println(t, string(newBytes))
	if err := t.Validate(); err != nil {
		return nil, err
	}
	tr.db.Lock()
	if err = tr.db.Pdb.Model(&models.Task{}).Where("id = ?", t.ID).Save(&oldObj).Error;
	err != nil {
		return nil, err
	}
	tr.db.Unlock()
	
	return t, nil
}