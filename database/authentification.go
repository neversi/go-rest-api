package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	// "gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
	
)

// DB here
type DB struct {
	db *sql.DB
}

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
			Addr: "localhost: 5404",
			DB: 0,
			Password: "",
	})
)

func (db *DB) findUser(login string) (*models.User, error) {
	return nil, nil
} 

func (db *DB) getUser(login string) (*models.User, error) {
	
	u, err := db.findUser(login)
	if err != nil {
		err = fmt.Errorf("The user with such login does not exist")
		return nil, err
	}
	
	return u, nil
}

// TokenDetails structure with info about access and refresh tokens
type TokenDetails struct {
	AToken 	string
	RToken 	string
	AUuid 	string
	RUuid 	string
	AExp	int64
	RExp	int64
}

// AccessToken ...
type AccessToken struct {
	AUuid string
	Userid string
}

// RefreshToken ...
type RefreshToken struct {
	RUuid string
	Userid string
}



// Login func validates login and password and returns the token
func Login(w http.ResponseWriter, r *http.Request) {
	var u = models.User{}
	bodyBytes, err := ioutil.ReadAll(r.Body);
	if err != nil {
		err = fmt.Errorf("Cannot read the error")
		return
	}

	err = json.Unmarshal(bodyBytes, &u)
	if err != nil {
		err = fmt.Errorf("Cannot read the error")
		return
	}

	db := DB{}
	checkUser, err := db.getUser(u.Login)
	if err != nil {
		err = fmt.Errorf("Cannot read the error")
		return
	}
	
	if valid, _ := CompareHash(checkUser.Password, u.Password); valid !=  true  {
		w.Write([]byte("StatusUnauthorized"));
		w.WriteHeader(http.StatusUnauthorized);
		return
	}
	token, err := CreateToken(u.ID)

	
	_ = token
}

// CreateToken creates an access and refresh tokens
func CreateToken(id uint) (*TokenDetails, error) {
	token := &TokenDetails{}
	var err error
	token.AExp = time.Now().Add(time.Minute * 15).Unix()
	token.AUuid = uuid.New().String()
	token.RExp = time.Now().Add(time.Hour * 24 * 7).Unix()
	token.RUuid = uuid.New().String()

	envSecret := "abdr_go_to_env"

	AClaims := jwt.MapClaims{}

	AClaims["authorized"] = true
	AClaims["a_id"] = token.AUuid
	AClaims["id"] = id
	AClaims["exp"] = token.AExp
	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, AClaims)

	token.AToken, err = atoken.SignedString(envSecret)
	if err != nil {
		err = fmt.Errorf("signing error")
		return nil, err
	}

	RClaims := jwt.MapClaims{}
	RClaims["authorized"] = true
	RClaims["r_id"] = token.RUuid
	RClaims["id"] = id
	RClaims["exp"] = token.RExp
	rtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, RClaims)
	
	token.RToken, err = rtoken.SignedString(envSecret)
	if err != nil {
		err = fmt.Errorf("signing error")
		return nil, err
	}

	return token, nil
}

//CreateAuth creates a valid token in cached db
func CreateAuth(id int64, token *TokenDetails) error {
	at := time.Unix(token.AExp, 0)
	rt := time.Unix(token.RExp, 0)
	err := rdb.Set(ctx, token.AUuid, strconv.Itoa(int(id)), at.Sub(time.Now())).Err()
	if err != nil {
		err = fmt.Errorf("error with creating the cached memory for access token")
		return err
	}

	err = rdb.Set(ctx, token.RUuid, strconv.Itoa(int(id)), rt.Sub(time.Now())).Err()

	if err != nil {
		err = fmt.Errorf("error with creatig the cached memory for refresh token")
		return err
	}

	return nil
}


// ExtractToken ...
func ExtractToken(r *http.Request) string {
	rawToken := r.Header.Get("Authorization")
	subArr := strings.Split(rawToken, " ")

	if (len(subArr) == 2) {
		return subArr[1]
	}

	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		}

		return []byte(os.Getenv("secret")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// TokenValid ...
func TokenValid(r *http.Request) (error) {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

// ExtractTokenData ... 
func ExtractTokenData(r *http.Request) (*AccessToken, error) {
	token, err := verifyToken(r)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("no claims")
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 0)
	if err != nil {
		return nil, err
	}
	accessUUID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["a_id"]), 10, 0)
	if err != nil {
		return nil, err
	}

	
	return &AccessToken{
		AUuid: strconv.Itoa(int(accessUUID)),
		Userid: strconv.Itoa(int(userID)),
	}, nil
}

// CheckAuth ...
func CheckAuth(authD *AccessToken) (int, error) {
	userID, err := rdb.Get(ctx, authD.AUuid).Result()
	if err != nil {
		return -1, nil
	}

	userIDint, _ := strconv.ParseInt(userID, 10, 0)
	return int(userIDint), nil
}

// DeleteAuth ...
func DeleteAuth(uuid string) (int64, error) {
	deleted, err := rdb.Del(ctx, uuid).Result()
	if err != nil {
		return -1, err
	}
	return deleted, nil
}

// Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {
	au, err := ExtractTokenData(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	delID, err := DeleteAuth(au.AUuid)

	if err != nil && delID != 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully logged out\n"))
}

// Refresh ... 
func Refresh(w http.ResponseWriter, r *http.Request) {
	tokens := map[string]string{}

	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	if err = json.Unmarshal(bodyBytes, &tokens); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	refreshToken := tokens["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		   return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	     })
	
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
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

		if deleted, ok := DeleteAuth(refreshUUID); deleted != 0 && ok != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		userID := mapClaims["id"].(string)
		userIDInt, _ := strconv.Atoi(userID)
		refreshedTokens, err := CreateToken(uint(userIDInt))
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
		}

		err = CreateAuth(int64(userIDInt), refreshedTokens)
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