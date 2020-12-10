package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/cache"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/service"
)

// TaskController handles the request of the user
type TaskController struct {
	taskService service.ITaskService
	postCache cache.PostCache
}

// NewTaskController creates TaskController
func NewTaskController(db *database.DataBase, rc *cache.RedisCache) *TaskController {
	return &TaskController{
		taskService: &service.TaskService{
			TaskRepository: database.NewTaskRepository(db),
		}, postCache: rc,
	}
}

// Create handles the request and creates
func (controller *TaskController) Create(w http.ResponseWriter, r *http.Request) {
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

	if err = controller.taskService.Create(task); err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}

	misc.JSONWrite(w, misc.WriteResponse(false, "Task created"), http.StatusCreated)
}

// Read retrieves certain datas from DB
func (controller *TaskController) Read(w http.ResponseWriter, r *http.Request) {
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
		task, err := controller.postCache.Get(userIDs) 
		if task == nil {
			task, err = controller.taskService.FindByID(uint(userID))

			if err != nil {
				misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
				return
			}
			tasks, err = controller.taskService.Read(task)
			controller.postCache.Set(userIDs, task)
		}
		
	} else {
		tasks, err = controller.taskService.Read(nil)
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
func (controller *TaskController) Delete(w http.ResponseWriter, r *http.Request) {
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
		
		task, err := controller.taskService.FindByID(uint(taskID))
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusOK)
			return
		}
		if uint(userID) == task.UserID {
			err = controller.taskService.Delete(task)
		}
	}

	if err != nil { 
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
		return
	}

	misc.JSONWrite(w, misc.WriteResponse(false, "Successfully deleted"), http.StatusOK)
}

// Update updates the tasks parameters in DB
func (controller *TaskController) Update(w http.ResponseWriter, r *http.Request) {
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
		task, err = controller.taskService.FindByID(uint(taskID))
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
	
	err = controller.taskService.Update(task)

	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}

	misc.JSONWrite(w, misc.WriteResponse(false, "Successfully updated"), http.StatusOK)
}

// GetByID ... 
func (controller *TaskController) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.ParseInt(vars["task_id"], 10, 0)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}
	task, err := controller.taskService.FindByID(uint(taskID))

	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusNotFound)
		return
	}
	
	misc.JSONWrite(w, misc.WriteResponse(false, task), http.StatusFound)
}