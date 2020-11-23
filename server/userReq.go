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

// UserReq handles the request of the user
type UserReq struct {
	api *APIServer
}

// NewUserReq creates UserReq
func NewUserReq(api *APIServer) *UserReq {
	return &UserReq{api: api}
}

// UserCreate creates the user in the database
func (ur *UserReq) UserCreate(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	bodyBytes, err := ioutil.ReadAll(r.Body)
	fmt.Println(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, &user)
	
	_, err = ur.api.db.User().Create(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UserRead responds with json format file where all users are written
func (ur *UserReq) UserRead(w http.ResponseWriter, r *http.Request) {
	users, err := ur.api.db.User().Read(nil)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	for _, user := range users {
		fmt.Fprintf(w, "User #%d %s %s\n", user.ID, user.Login, user.Email)
	}
	w.WriteHeader(http.StatusOK)
}

// UserUpdate updates the info about user
func (ur *UserReq) UserUpdate(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	user := new (models.User)

	err = json.Unmarshal(bodyBytes, &user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err = ur.api.db.User().Update(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

// UserDelete deletes the user account
func (ur *UserReq) UserDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user := new (models.User)

	userID, _ := strconv.ParseInt(vars["id"], 10, 0)
	
	ur.api.db.Pdb.Raw("SELECT login FROM users WHERE id = ?", uint(userID)).Scan(&user.Login)
	
	if err := ur.api.db.User().Delete(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	
}