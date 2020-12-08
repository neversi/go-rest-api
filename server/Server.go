package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/configs"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/middleware"
)

// APIServer is server which will have all flows and so on
type APIServer struct {
	config 	       *configs.Server
	userController *UserController 		// User
	taskController *TaskController 
}

// New initialize the server
func New()(*APIServer, error) {
	var configs *configs.Server
	db := database.New()
	db.OpenDataBase()

	return &APIServer{
		taskController: NewTaskController(db),
		userController: NewUserController(db),
	}, nil
}

// Start starts the server
func (api *APIServer) Start() error {
	router := mux.NewRouter()
	// subRoutes
	// simpler as you can
	// Redis - catalog
	router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
		w.WriteHeader(http.StatusOK)
	})
	userRouter := router.PathPrefix("/users").Subrouter()
	taskRouter := router.PathPrefix("/tasks").Subrouter()
	userRouter.HandleFunc("/", api.userController.Read).Methods("GET")
	router.HandleFunc("/login", api.userController.Login).Methods("POST")
	// router.HandleFunc("/register", api.userController.Register).Methods("POST")
	taskRouter.HandleFunc("/", api.taskController.Create).Methods("POST")
	userRouter.HandleFunc("/", api.userController.Create).Methods("POST")
	userRouter.HandleFunc("/{id}", api.userController.Delete).Methods("DELETE")
	userRouter.HandleFunc("/{id}", api.userController.Update).Methods("PUT")
	

	
	// router.HandleFunc("/refresh", api.userController.Refresh).Methods("GET")
	userRouter.HandleFunc("/{id}/tasks",api.taskController.Read).Methods("GET")
	userRouter.HandleFunc("/{id}/tasks/{task_id}", api.taskController.Delete).Methods("DELETE")
	userRouter.HandleFunc("/{id}/tasks/{task_id}", api.taskController.Update).Methods("PUT")
	router.HandleFunc("/tasks", api.taskController.Read).Methods("GET")

	router.Use(middleware.JSONDataCheck)
	userRouter.Use(middleware.IsAuthorized)
	userRouter.Use(middleware.AuthorizationUser)
	
	return http.ListenAndServe(api.config.Port, router)
}

// Greeting function
func Greeting(w http.ResponseWriter, r *http.Request) {
}
