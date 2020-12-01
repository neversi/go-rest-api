package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/service"
)

// TaskController handles the request of the user
type TaskController struct {
	taskService service.ITaskService
}

// NewTaskController creates TaskController
func NewTaskController(db *database.DataBase) *TaskController {
	return &TaskController{
		taskService: &service.TaskService{
			TaskRepository: database.NewTaskRepository(db),
		},
	}
}

// Create handles the request and creates
func (tr *TaskController) Create(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	task := new(models.Task)

	err = json.Unmarshal(bodyBytes, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err = task.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = tr.taskService.Create(task); err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}

	misc.JSONWrite(w, misc.WriteResponse(false, "Success"), http.StatusCreated)
}

// Read retrieves certain datas from DB
func (tr *TaskController) Read(w http.ResponseWriter, r *http.Request) {
	var tasks []*models.Task
	var err error
	vars := mux.Vars(r)
	userIDs, ok := vars["id"]
	if (ok) {
		userID, err := strconv.ParseInt(userIDs, 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		task := &models.Task{UserID: uint(userID)}
		err = tr.taskService.Read(task)
		
	} else {
		err = tr.taskService.Read(nil)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		fmt.Fprintf(w, "User #%d %s %s\n", task.ID, task.Title, task.Author)
	}
}

// Delete deletes task from DB
func (tr *TaskController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	taskIDs, ok := vars["task_id"]
	if ok {
		taskID, err := strconv.ParseInt(taskIDs, 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		task := &models.Task{ID: uint(taskID)}
		err = tr.taskService.Delete(task)
	}
	// title implement!
	if err != nil { 
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Update updates the tasks parameters in DB
func (tr *TaskController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	var task *models.Task
	taskIDs, ok := vars["task_id"]
	userIDs, _ := vars["id"]
	if ok {
		taskID, err := strconv.ParseInt(taskIDs, 10, 0)
		userID, err := strconv.ParseInt(userIDs, 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		task = &models.Task{ID: uint(taskID), UserID: uint(userID)}
	}
	
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.Unmarshal(bodyBytes, &task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	fmt.Println(task)
	err = tr.taskService.Update(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
}