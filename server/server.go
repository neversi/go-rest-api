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
	userController *UserController 		
	taskController *TaskController 
}

// New initialize the server
func New(conf *configs.Server) (*APIServer, error) {
	db := database.New()
	db.OpenDataBase(conf.DB)
	cacheHost := conf.Cache.Host + ":" + conf.Cache.Port
	cache := cache.NewRedisCache(cacheHost, conf.Cache.DB, time.Duration(conf.Cache.ExpDuration))
	fmt.Print()
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

	userRouter := router.PathPrefix("v1/users").Subrouter()
	taskRouter := router.PathPrefix("v1/tasks").Subrouter()
	adminRouter := router.PathPrefix("v1/admin").Subrouter()

	router.HandleFunc("/v1/login", api.userController.Login).Methods("POST")
	router.HandleFunc("/v1/register", api.userController.Create).Methods("POST")
	router.HandleFunc("/v1/refresh", Refresh).Methods("GET")

	adminRouter.HandleFunc("/users", api.userController.Read).Methods("GET")
	adminRouter.HandleFunc("/users", api.userController.Create).Methods("POST")
	adminRouter.HandleFunc("/users/{id:[0-9]+}", api.userController.Delete).Methods("DELETE")
	adminRouter.HandleFunc("/users/{id:[0-9]+}", api.userController.Update).Methods("PUT")
	adminRouter.HandleFunc("/tasks/{task_id:[0-9]+}", api.taskController.Delete).Methods("DELETE")
	adminRouter.HandleFunc("/tasks", api.taskController.Read).Methods("GET")

	userRouter.HandleFunc("/{id:[0-9]+}", api.userController.Delete).Methods("DELETE")
	userRouter.HandleFunc("/{id:[0-9]+}", api.userController.Update).Methods("PUT", "POST")
	userRouter.HandleFunc("/{id:[0-9]+}/tasks",api.taskController.Create).Methods("POST")
	userRouter.HandleFunc("/{id:[0-9]+}/tasks",api.taskController.Read).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}/tasks/{task_id:[0-9]+}", api.taskController.Delete).Methods("DELETE")
	userRouter.HandleFunc("/{id:[0-9]+}/tasks/{task_id:[0-9]+}", api.taskController.Update).Methods("PUT", "POST")
	userRouter.HandleFunc("/{id:[0-9]+}/tasks/{task_id:[0-9]+}", api.taskController.GetByID).Methods("GET")
	
	taskRouter.HandleFunc("/", api.taskController.Create).Methods("POST")
	
	
	router.Use(middleware.JSONDataCheck)
	router.Use(middleware.LoggerHandler)

	userRouter.Use(middleware.IsAuthenticated)
	userRouter.Use(middleware.AuthorizationUser)

	return http.ListenAndServe(fmt.Sprintf(":%s", api.config.Port), router)
	
}