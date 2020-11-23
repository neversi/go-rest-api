package server

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