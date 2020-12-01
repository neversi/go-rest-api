package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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
			database.NewUserRepository(db),
		},
	}
}

// Create creates the user in the database
func (ur *UserController) Create(w http.ResponseWriter, r *http.Request)  {
	if r.Header.Get("Content-Type") != "application/json" {
		misc.JSONWrite(w, misc.WriteResponse(true, "Unsupported Media Type: need application/json"), http.StatusUnsupportedMediaType)
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
		return
	}

	newUser := new(models.User)

	_ = json.Unmarshal(bodyBytes, &newUser)

	err = ur.userService.Create(newUser)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
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
	if (len(bodyBytes) != 0) {
		isEmpty = !isEmpty
	}
	users := make([]*database.UserDTO, 0)
	if isEmpty == false {
		var user *database.UserDTO
		err = json.Unmarshal(bodyBytes, &user)
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		}
		users, err := ur.userService.FindByID()
	} else {
		users, err := ur.userService.Read(nil)
	}
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	misc.JSONWrite(w, misc.WriteResponse(false, users), http.StatusOK)

}

// UserUpdate updates the info about user
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
	
	
	// var user *models.User
	// err = json.Unmarshal(bodyBytes, &user)
	// users := v.([]*models.User)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	userIDi, err := strconv.ParseInt(vars["id"], 10, 0)
	userDTO.ID = uint(userIDi)
	// for _, user := range users {
		err = ur.userService.Update(userDTO)
		if err != nil {	
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	// }
	
	w.WriteHeader(http.StatusOK)
}

// Delete deletes the user account
func (ur *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		misc.JSONWrite(w, misc.WriteResponse(true, "Unsupported Media Type: need application/json"), http.StatusUnsupportedMediaType)
		return
	}
	vars := mux.Vars(r)

	userDTO := new(database.UserDTO)
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
	// contentType := r.Header.Get("Content-Type")
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

	users, err := ur.userService.FindByLogin(u.Login)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusNotFound)
		return
	}

	if users == nil {
		misc.JSONWrite(w, "There is no such user", http.StatusNotFound)
		return
	}

	checkUser := users[0]
	if valid, _ := misc.CompareHash(checkUser.Password, u.Password); valid !=  true  {
		misc.JSONWrite(w, "Wrong password", http.StatusUnauthorized)
		return
	}
	token, err := auth.CreateToken(u.ID)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}
	if err = auth.CreateAuth(checkUser.ID, token); err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}

	signedTokens := map[string]string{"accessToken":token.AToken, "refreshToken":token.RToken}
	jsonBytes, err := json.Marshal(signedTokens)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return		
	}
	w.Write(jsonBytes)
	w.WriteHeader(http.StatusCreated)
}

// Register registers the node
func (ur *UserController) Register(w http.ResponseWriter, r *http.Request) {
	ur.Create(w, r)
	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
	return

}

// Logout ...
func (ur *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	au, err := auth.ExtractTokenData(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	delID, err := auth.DeleteAuth(ur.userRepository.DB, au.AUuid)

	if err != nil && delID != 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Successfully logged out\n"))
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
		w.WriteHeader(http.StatusUnprocessableEntity)
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims); 
	if ok && token.Valid {
		refreshUUID, ok := mapClaims["r_id"].(string)
		if !ok {
			w.WriteHeader(http.StatusExpectationFailed)
		}

		if deleted, err := auth.DeleteAuth(refreshUUID); deleted != 0 && err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		userID := mapClaims["id"].(string)
		userIDInt, _ := strconv.Atoi(userID)
		refreshedTokens, err := auth.CreateToken(uint(userIDInt))
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
		}

		err = auth.CreateAuth(uint(userIDInt), refreshedTokens)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cannot authorize token")
			return
		}

		tokens := map[string]string{
			"access_token": refreshedTokens.AUuid,
			"refresh_token": refreshedTokens.RUuid,
		}
		bodyBytes, err = json.Marshal(tokens)
		w.Write(bodyBytes)
		w.WriteHeader(http.StatusCreated)
		return
	}
	
	w.WriteHeader(http.StatusUnauthorized)	
}