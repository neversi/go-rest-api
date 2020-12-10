package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/database"
)

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
	URole string
	AUuid string
	Userid uint64
}

// CreateToken creates an access and refresh tokens
func CreateToken(id uint, role string) (*TokenDetails, error) {
	token := &TokenDetails{}
	var err error
	token.AExp = time.Now().Add(time.Minute * 15).Unix()
	token.AUuid = uuid.New().String()
	token.RExp = time.Now().Add(time.Hour * 24).Unix()
	token.RUuid = uuid.New().String()


	AClaims := jwt.MapClaims{}

	AClaims["authorized"] = true
	AClaims["role"] = role
	AClaims["a_id"] = token.AUuid
	AClaims["id"] = id
	AClaims["exp"] = token.AExp
	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, AClaims)

	token.AToken, err = atoken.SignedString([]byte(os.Getenv("TokenPass")))
	if err != nil {
		return nil, err
	}

	RClaims := jwt.MapClaims{}
	RClaims["authorized"] = true
	RClaims["role"] = role
	RClaims["r_id"] = token.RUuid
	RClaims["id"] = id
	RClaims["exp"] = token.RExp
	rtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, RClaims)
	
	token.RToken, err = rtoken.SignedString([]byte(os.Getenv("TokenPass")))
	if err != nil {
		return nil, err
	}

	return token, nil
}

//CreateAuth creates a valid token in cached db
func CreateAuth(db *database.DataBase, id uint, token *TokenDetails) error {
	rt := time.Unix(token.RExp, 0)
	at := time.Unix(token.AExp, 0)

	err := db.Rdb.Set(db.Ctx, token.AUuid, strconv.Itoa(int(id)), at.Sub(time.Now())).Err()
	if err != nil {
		return err
	}

	err = db.Rdb.Set(db.Ctx, token.RUuid, strconv.Itoa(int(id)), rt.Sub(time.Now())).Err()
	if err != nil {
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
	} else if len(subArr) == 1 {
		return subArr[0]
	}

	return ""
}

// VerifyToken verifies if token is exist or not
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		}

		return []byte(os.Getenv("TokenPass")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// TokenValid validity of token
func TokenValid(r *http.Request) (error) {
	token, err := VerifyToken(r)
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
	token, err := VerifyToken(r)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("no claims")
	}

	userID := claims["id"]
	accessUUID := claims["a_id"].(string)

	userRole := claims["role"].(string)
	
	return &AccessToken{
		URole: userRole,
		AUuid: accessUUID,
		Userid: uint64(userID.(float64)),
	}, nil
}

// CheckAuth ...
func CheckAuth(db *database.DataBase, authD *AccessToken) (int, error) {
	userID, err := db.Rdb.Get(db.Ctx, authD.AUuid).Result()
	if err != nil {
		return -1, nil
	}

	userIDint, _ := strconv.ParseInt(userID, 10, 0)
	return int(userIDint), nil
}


// DeleteAuth ...
func DeleteAuth(db *database.DataBase, uuid string) (int64, error) {
	deleted, err := db.Rdb.Del(db.Ctx, uuid).Result()
	if err != nil {
		return -1, err
	}
	return deleted, nil
}