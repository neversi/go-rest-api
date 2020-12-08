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
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}
	task := new(models.Task)

	err = json.Unmarshal(bodyBytes, &task)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
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
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
			return
		}
		task, err := tr.taskService.FindByID(uint(userID))
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
			return
		}
		tasks, err = tr.taskService.Read(task)
		
	} else {
		tasks, err = tr.taskService.Read(nil)
	}
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
		return
	}

	if tasks == nil {
		misc.JSONWrite(w, misc.WriteResponse(false, "There is no tasks"), http.StatusOK)
		return
	}

	misc.JSONWrite(w, misc.WriteResponse(false, tasks), http.StatusOK)
}

// Delete deletes task from DB
func (tr *TaskController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := fmt.Errorf("The id of task was not declared")
	taskIDs, ok := vars["task_id"]
	if ok {
		taskID, err := strconv.ParseInt(taskIDs, 10, 0)
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
			return
		}
		userID, err := strconv.ParseInt(vars["id"], 10, 0)
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
			return
		}
		
		task, err := tr.taskService.FindByID(uint(taskID))
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusOK)
			return
		}
		if uint(userID) == task.UserID {
			err = tr.taskService.Delete(task)
		}
	}

	if err != nil { 
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
		return
	}

	misc.JSONWrite(w, misc.WriteResponse(false, "Successfully deleted"), http.StatusOK)
}

// Update updates the tasks parameters in DB
func (tr *TaskController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	var task *models.Task
	taskIDs, ok := vars["task_id"]
	if ok {
		taskID, err := strconv.ParseInt(taskIDs, 10, 0)
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
			return
		}
		userID, err := strconv.ParseInt(vars["id"], 10, 0)
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
			return
		}
		task, err = tr.taskService.FindByID(uint(taskID))
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
			return
		}
		if task.UserID != uint(userID) {
			misc.JSONWrite(w, misc.WriteResponse(true, "Not authorized user"), http.StatusUnprocessableEntity)
			return
		}
	}
	
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}

	err = json.Unmarshal(bodyBytes, &task)

	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}

	taskID, _ := strconv.ParseInt(taskIDs, 10, 0)
	userID, _ := strconv.ParseInt(vars["id"], 10, 0)

	task.ID = uint(taskID)
	task.UserID = uint(userID)
	
	err = tr.taskService.Update(task)

	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}

	misc.JSONWrite(w, misc.WriteResponse(false, "Successfully updated"), http.StatusOK)
}