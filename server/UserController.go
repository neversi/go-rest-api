package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/auth"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/service"
)

// UserController handles the request of the user
type UserController struct {
	userService service.IUserService
}


// NewUserController creates UserController
func NewUserController(db *database.DataBase) *UserController {
	return &UserController{
		userService: &service.UserService{
			UserRepository: database.NewUserRepository(db),
		},
	}
}

// Create creates the user in the database
func (ur *UserController) Create(w http.ResponseWriter, r *http.Request)  {
	
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
		return
	}

	newUser := new(models.User)

	_ = json.Unmarshal(bodyBytes, &newUser)

	err = ur.userService.Create(newUser)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
		return
	}
	misc.JSONWrite(w, misc.WriteResponse(false, "User created"), http.StatusCreated)
}

// UserRead responds with json format file where all users are written
func (ur *UserController) Read(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		misc.JSONWrite(w, misc.WriteResponse(true, "Unsupported Media Type: need application/json"), http.StatusUnsupportedMediaType)
		return
	}
	isEmpty := true
	
	bodyBytes, err := json.Marshal(r.Body)
	if err != nil { 
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}
	if (string(bodyBytes) != "{}") {
		isEmpty = false
	}
	users := make([]*database.UserDTO, 0)
	if isEmpty == false {
		var user *database.UserDTO
		err = json.Unmarshal(bodyBytes, &user)
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		}
		user, err = ur.userService.FindByID(user.ID);
		users = append(users, user);
	} else {
		users, err = ur.userService.Read(nil)
	}
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	misc.JSONWrite(w, misc.WriteResponse(false, users), http.StatusOK)

}

// Update updates the info about user
func (ur *UserController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		misc.JSONWrite(w, misc.WriteResponse(true, "Unsupported Media Type: need application/json"), http.StatusUnsupportedMediaType)
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		misc.JSONWrite(w,misc.WriteResponse(true, err.Error()),http.StatusUnprocessableEntity)
		return
	}
	vars := mux.Vars(r)
	

	userDTO := new(database.UserDTO)
	
	err = json.Unmarshal(bodyBytes, &userDTO)
	

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	userIDi, err := strconv.ParseInt(vars["id"], 10, 0)
	userDTO.ID = uint(userIDi)
	err = ur.userService.Update(userDTO)
	if err != nil {	
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}
	
	misc.JSONWrite(w, misc.WriteResponse(false, "Successfully updated"), http.StatusOK)
}

// Delete deletes the user account
func (ur *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		misc.JSONWrite(w, misc.WriteResponse(true, "Unsupported Media Type: need application/json"), http.StatusUnsupportedMediaType)
		return
	}
	vars := mux.Vars(r)
	userDTO := new(database.UserDTO)
	bodyBytes, err := ioutil.ReadAll(r.Body);
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &userDTO)
	userID, _ := strconv.ParseInt(vars["id"], 10, 0)
	userDTO.ID = uint(userID)
	if err := ur.userService.Delete(userDTO); err != nil {
		misc.JSONWrite(w,misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}
	misc.JSONWrite(w, misc.WriteResponse(false, "Success"), 204)
}

// Login authentificate the user by checking and giving the token
func (ur *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		misc.JSONWrite(w, misc.WriteResponse(true, "Unsupported Media Type: need application/json"), http.StatusUnsupportedMediaType)
		return
	}
	u := new(models.User)

	bodyBytes, err := ioutil.ReadAll(r.Body);
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, &u)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}

	err = ur.userService.CheckUser(u.Login, u.Password)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusNotFound)
		return
	}

	user, err := ur.userService.FindByLogin(u.Login);

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
	}

	signedTokens := map[string]string{"accessToken":token.AToken, "refreshToken":token.RToken}

	misc.JSONWrite(w, misc.WriteResponse(false, signedTokens), http.StatusCreated)

}

// Register registers the node
func (ur *UserController) Register(w http.ResponseWriter, r *http.Request) {
	ur.Create(w, r)
	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
	return

}

// Logout ...
func (ur *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractTokenData(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	misc.JSONWrite(w, misc.WriteResponse(false, "Successfully logged out"), http.StatusAccepted)
}



// Refresh ... 
func Refresh(w http.ResponseWriter, r *http.Request) {
	tokens := map[string]string{}

	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err = json.Unmarshal(bodyBytes, &tokens); err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}

	refreshToken := tokens["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		   return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("should change"), nil
	     })
	
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnauthorized)
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims); 
	if ok && token.Valid {
		// refreshUUID, ok := mapClaims["r_id"].(string)
		// if !ok {
		// 	w.WriteHeader(http.StatusExpectationFailed)
		// }

		userID := mapClaims["id"].(string)
		userIDInt, _ := strconv.Atoi(userID)
		refreshedTokens, err := auth.CreateToken(uint(userIDInt))
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
		}

		tokens := map[string]string{
			"access_token": refreshedTokens.AUuid,
			"refresh_token": refreshedTokens.RUuid,
		}
		misc.JSONWrite(w, misc.WriteResponse(false, tokens), http.StatusCreated)
		return
	}
	
	w.WriteHeader(http.StatusUnauthorized)	
}