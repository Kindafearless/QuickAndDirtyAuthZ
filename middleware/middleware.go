package middleware

import (
	"net/http"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"github.com/gorilla/context"
	"QuickAndDirtyAuthZ/authorization"
	"github.com/gorilla/mux"
)

type Token struct {
	user string
	role string
	exp string
}

const secretKey = "igorigorigorigorigor"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		fmt.Printf("Request: %v\n", req.RequestURI)

		authHeader := req.Header.Get("Authorization")

		// if the header exists, setup context
		if authHeader == "" || strings.Contains(authHeader, "Basic") {
			http.Error(rw, "Request requires token", http.StatusUnauthorized)
			return
		}

		// Decode the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := decodeToken(tokenString)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusUnauthorized)
			return
		}

		// Map claims. Add role to context, check if endp
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			context.Set(req, "role", claims["role"])
			context.Set(req, "user", claims["user"])

			authorized := authorization.ConfirmRoleMapping(claims["role"].(string), authorization.DecodeRolePath(req.URL.Path, mux.Vars(req)))
			if !authorized {
				http.Error(rw, "forbidden", http.StatusUnauthorized)
				return
			}

		} else {
			http.Error(rw, err.Error(), http.StatusUnauthorized)
			return
		}

		// Pass on to the next handler
		next.ServeHTTP(rw, req)
	})
}

func decodeToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(secretKey), nil
}