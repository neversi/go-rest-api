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
	roleService service.IRoleService
}


// NewUserController creates UserController
func NewUserController(db *database.DataBase) *UserController {
	return &UserController{
		userService: &service.UserService{
			UserRepository: database.NewUserRepository(db),
		},
		roleService: &service.RoleService{
			RoleRepository: database.NewRoleRepository(db),
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
	newRole := new(models.Role)

	err = json.Unmarshal(bodyBytes, &newUser)
	
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}
	
	err = json.Unmarshal(bodyBytes, &newRole)
	if newRole.Role == "" {
		newRole.Role = "user"
	}
	
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}
	
	err = ur.userService.Create(newUser)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusBadRequest)
		return
	}

	newUser, err = ur.userService.FindByLogin(newUser.Login)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}

	err = ur.roleService.SetUserRole(newUser.ID, newRole.Role)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}
	
	misc.JSONWrite(w, misc.WriteResponse(false, "User created"), http.StatusCreated)
}

// UserRead responds with json format file where all users are written
func (ur *UserController) Read(w http.ResponseWriter, r *http.Request) {
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
	
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		misc.JSONWrite(w,misc.WriteResponse(true, err.Error()),http.StatusUnprocessableEntity)
		return
	}
	vars := mux.Vars(r)
	

	userDTO := new(database.UserDTO)
	newRole := new(models.Role)

	err = json.Unmarshal(bodyBytes, &userDTO)

	if err != nil {
		misc.JSONWrite(w,misc.WriteResponse(true, err.Error()),http.StatusUnprocessableEntity)
		return
	}
	
	err = json.Unmarshal(bodyBytes, &newRole) 
	
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnprocessableEntity)
		return
	}
	
	userIDi, err := strconv.ParseInt(vars["id"], 10, 0)
	userDTO.ID = uint(userIDi)
	err = ur.userService.Update(userDTO)
	if err != nil {	
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}

	usr, err := ur.userService.FindByLogin(userDTO.Login)
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}

	err = ur.roleService.SetUserRole(usr.ID, newRole.Role)
	
	if err != nil {
		misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
		return
	}
	misc.JSONWrite(w, misc.WriteResponse(false, "Successfully updated"), http.StatusOK)
}

// Delete deletes the user account
func (ur *UserController) Delete(w http.ResponseWriter, r *http.Request) {
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

	user, err := ur.userService.FindByLogin(u.Login)
	role, err := ur.roleService.FindByUserID(user.ID)
	token, err := auth.CreateToken(user.ID, role)
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
		_, ok := mapClaims["r_id"].(string)
		if !ok {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusExpectationFailed)
			return
		}

		userID := mapClaims["id"].(string)
		userIDInt, _ := strconv.Atoi(userID)
		refreshedTokens, err := auth.CreateToken(uint(userIDInt), mapClaims["role"].(string))
		if err != nil {
			misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusInternalServerError)
			return
		}

		tokens := map[string]string{
			"access_token": refreshedTokens.AUuid,
			"refresh_token": refreshedTokens.RUuid,
		}
		misc.JSONWrite(w, misc.WriteResponse(false, tokens), http.StatusCreated)
		return
	}
	
	misc.JSONWrite(w, misc.WriteResponse(true, err.Error()), http.StatusUnauthorized)
	return
}