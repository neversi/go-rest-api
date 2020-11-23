package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
)

// APIServer is server which will have all flows and so on
type APIServer struct {
	config string
	db *database.DataBase
	userReq *UserReq
	taskReq *TaskReq
}

// New initialize the server
func New() (*APIServer, error) {
	db := database.New()
	db.OpenDataBase()
	newAPI := &APIServer{}
	newAPI.db = db
	newAPI.userReq = NewUserReq(newAPI)
	newAPI.taskReq = NewTaskReq(newAPI)
	
	return newAPI, nil
}

// Start starts the server
func (api *APIServer) Start() error {
	router := mux.NewRouter()
	
	router.HandleFunc("/", Example)
	router.HandleFunc("/users", api.userReq.UserRead).Methods("GET")
	router.HandleFunc("/users/create", api.userReq.UserCreate).Methods("POST")
	router.HandleFunc("/users/delete/{id}", api.userReq.UserDelete).Methods("DELETE")
	router.HandleFunc("/users/update/{id}", api.userReq.UserUpdate).Methods("PUT")
	return http.ListenAndServe(":8080", router)
}

// Example function
func Example(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
	w.WriteHeader(http.StatusOK)
}