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
	userController *UserController // User
	taskController *TaskController 
}

// New initialize the server
func New()(*APIServer, error) {
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
	router.HandleFunc("/", Greeting)
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
	userRouter.HandleFunc("/{id}/tasks/{task_id}", api.taskController.Delete).Methods("PUT")
	router.HandleFunc("/tasks", api.taskController.Read).Methods("GET")

	
	userRouter.Use(isAuthorized)
	return http.ListenAndServe(":8080", router)
}

// Greeting function
func Greeting(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
	w.WriteHeader(http.StatusOK)
}

// !!! Move to middleware dirallows modeling many natural phenomena
func isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		if r.Header["Authorization"] != nil {
			if err := auth.TokenValid(r); err != nil {
				http.Error(w, "Token expired", http.StatusForbidden)
			} else {
				acessToken, err := auth.ExtractTokenData(r)
				if err != nil {
					http.Error(w, "Problem with Token", http.StatusUnauthorized)
					return
				}
				_ = acessToken
				next.ServeHTTP(w, r)
				return
			}
			
		} else {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
		}
	})
}