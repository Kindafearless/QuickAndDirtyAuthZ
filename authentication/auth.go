package authentication

import (
	"net/http"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"time"
	"encoding/base64"
	"strings"
	"AuthZ/database"
)

type DatabaseConnector struct {
	DB           *database.InMemory
	DatabaseName string
}

type DataObject struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

const secretKey = "igorigorigorigorigro"

func (dc *DatabaseConnector) Login(rw http.ResponseWriter, req *http.Request) {

	authHeader := req.Header.Get("Authorization")
	basic := strings.TrimPrefix(authHeader, "Basic ")

	decoded, err := base64.StdEncoding.DecodeString(basic)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	// retrieve the username and password from the decoded basic auth string
	pair := strings.SplitN(string(decoded), ":", 2)

	// ignore the fact that we aren't checking passwords. This is an AuthZ example, not an AuthN example :)

	// get user object from database (if they exist)
	user, err := dc.DB.GetFromDatabase(dc.DatabaseName, pair[0])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := makeJWT(user.Key, user.Value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	resp, err := json.Marshal(token)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func makeJWT(user string, role string) (string, error) {

	// In a real authN/Z setup, your identity provider would be granting you a token with more information than this.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"role": role,
		"exp": time.Now().Add(time.Minute * 30).String(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil { return "", err }

	return tokenString, nil
}

func (dc *DatabaseConnector) List(rw http.ResponseWriter, req *http.Request) {

	transactions, err := dc.DB.ListDatabase(dc.DatabaseName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(transactions)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}