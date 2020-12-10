package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/cache"
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
func New(conf *configs.Server) (*APIServer, error) {
	db := database.New()
	db.OpenDataBase(conf.DB)
	cacheHost := conf.Cache.Host + ":" + conf.Cache.Port
	cache := cache.NewRedisCache(cacheHost, conf.Cache.DB, time.Duration(conf.Cache.ExpDuration))
	
	return &APIServer{
		config: conf,
		taskController: NewTaskController(db, cache),
		userController: NewUserController(db),
	}, nil
}

// Start starts the server
func (api *APIServer) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
		w.WriteHeader(http.StatusOK)
	})
	userRouter := router.PathPrefix("/users").Subrouter()
	taskRouter := router.PathPrefix("/tasks").Subrouter()

	router.HandleFunc("/login", api.userController.Login).Methods("POST")
	router.HandleFunc("/register", api.userController.Create).Methods("POST")

	router.HandleFunc("/users", api.userController.Read).Methods("GET")
	router.HandleFunc("/users", api.userController.Create).Methods("POST")
	userRouter.HandleFunc("/{id:[0-9]+}", api.userController.Delete).Methods("DELETE")
	userRouter.HandleFunc("/{id:[0-9]+}", api.userController.Update).Methods("PUT")
	
	taskRouter.HandleFunc("/", api.taskController.Create).Methods("POST")
	
	// router.HandleFunc("/refresh", api.userController.Refresh).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}/tasks",api.taskController.Create).Methods("POST")
	userRouter.HandleFunc("/{id:[0-9]+}/tasks",api.taskController.Read).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}/tasks/{task_id:[0-9]+}", api.taskController.Delete).Methods("DELETE")
	userRouter.HandleFunc("/{id:[0-9]+}/tasks/{task_id:[0-9]+}", api.taskController.Update).Methods("PUT")
	router.HandleFunc("/tasks", api.taskController.Read).Methods("GET")

	userRouter.Use(middleware.IsAuthenticated)
	router.Use(middleware.JSONDataCheck)
	// router.Use(middleware.LoggerHandler)
	userRouter.Use(middleware.AuthorizationUser)

	return http.ListenAndServe(fmt.Sprintf(":%s", api.config.Port), router)
}


