package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/auth"
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
	
	router. HandleFunc("/users", api.userReq.UserRead).Methods("GET")
	router.HandleFunc("/login", api.userReq.Login).Methods("POST")

	router.HandleFunc("/users/create", api.userReq.UserCreate).Methods("POST")
	router.HandleFunc("/users/delete/{id}", api.userReq.UserDelete).Methods("DELETE")
	router.HandleFunc("/users/update/{id}", api.userReq.UserUpdate).Methods("PUT")
	router.HandleFunc("/tasks/create", api.taskReq.TaskCreate).Methods("POST")
	
	router.HandleFunc("/users/{id}/tasks",api.taskReq.TaskRead).Methods("GET")
	router.HandleFunc("/users/{id}/tasks/delete/{task_id}", api.taskReq.TaskDelete).Methods("DELETE")
	router.HandleFunc("/users/{id}/tasks/delete/{title}", api.taskReq.TaskDelete).Methods("DELETE")

	router.HandleFunc("/tasks", api.taskReq.TaskRead).Methods("GET")
	router.HandleFunc("/users/{id}/tasks/update/{task_id}", api.taskReq.TaskUpdate).Methods("PUT")
	// router.Use(isAuthorized)
	return http.ListenAndServe(":8080", router)
}

// Example function
func Example(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
	w.WriteHeader(http.StatusOK)
}

func isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/login": 
			next.ServeHTTP(w, r)
			return
		case "/users":
			next.ServeHTTP(w, r)
			return
		}
		if r.Header["Authorization"] != nil {
			token, err := auth.VerifyToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}

			if token.Valid == true {
				next.ServeHTTP(w, r)
				return 
			}
			
			http.Error(w, "Token is not valid", http.StatusForbidden)
		} else {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
		}
	})
}