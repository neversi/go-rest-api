package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/auth"
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
	var users *models.User
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.Unmarshal(bodyBytes, &users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	
	// for _, user := range users {
		_, err = ur.api.db.User().Create(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	// }
	
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
		http.Error(w,err.Error(),http.StatusUnprocessableEntity)
		return
	}

	var user *models.User
	err = json.Unmarshal(bodyBytes, &user)
	// users := v.([]*models.User)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	// for _, user := range users {
		user, err = ur.api.db.User().Update(user)
		if err != nil {	
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	// }
	
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

// Login authentificate the user by checking and giving the token
func (ur *UserReq) Login(w http.ResponseWriter, r *http.Request) {
	var u = models.User{}
	bodyBytes, err := ioutil.ReadAll(r.Body);
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = json.Unmarshal(bodyBytes, &u)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	users, err := ur.api.db.User().Read(&u)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	checkUser := users[0]
	fmt.Println(checkUser.Login, checkUser.Password)
	if valid, _ := database.CompareHash(checkUser.Password, u.Password); valid !=  true  {
		w.Write([]byte("StatusUnauthorized"));
		w.WriteHeader(http.StatusUnauthorized);
		return
	}
	token, err := auth.CreateToken(u.ID)
	if err != nil {
		w.Write([]byte("Error while creation of token"));
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// if err = auth.CreateAuth(checkUser.ID, token); err != nil {
	// 	w.Write([]byte("Error while auth of token"));
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	signedTokens := map[string]string{"accessToken":token.AToken, "refreshToken":token.RToken}
	jsonBytes, err := json.Marshal(signedTokens)
	if err != nil {
		
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
	w.WriteHeader(http.StatusOK)
}