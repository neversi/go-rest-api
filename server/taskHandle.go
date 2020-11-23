package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
)

// import "gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"

// TaskReq handles the request of the user
type TaskReq struct {
	api *APIServer
}

// NewTaskReq creates TaskReq
func NewTaskReq(api *APIServer) *TaskReq {
	return &TaskReq{
		api: api,
	}
}

// TaskCreate creates the task and imports it in DB
func (tr *TaskReq) TaskCreate(w http.ResponseWriter, r *http.Request) {
	// var tasks []*models.Task
	var task *models.Task
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

	// for _, task := range tasks {
		task, err = tr.api.db.Task().Create(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	// }
	w.WriteHeader(http.StatusOK)
}

// TaskRead retrieves certain datas from DB
func (tr *TaskReq) TaskRead(w http.ResponseWriter, r *http.Request) {
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
		task := &models.Task{UserRefer: uint(userID)}
		tasks, err = tr.api.db.Task().Read(task)
		
	} else {
		tasks, err = tr.api.db.Task().Read(nil)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		fmt.Fprintf(w, "User #%d %s %s\n", task.ID, task.Title, task.Author)
	}
}

// TaskDelete deletes task from DB
func (tr *TaskReq) TaskDelete(w http.ResponseWriter, r *http.Request) {
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
		err = tr.api.db.Task().Delete(task)
	}
	// title implement!
	if err != nil { 
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TaskUpdate updates the tasks parameters in DB
func (tr *TaskReq) TaskUpdate(w http.ResponseWriter, r *http.Request) {
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
		task = &models.Task{ID: uint(taskID), UserRefer: uint(userID)}
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
	task, err = tr.api.db.Task().Update(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
}