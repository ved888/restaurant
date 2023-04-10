package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"restaurant/common"
	"strings"
)

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tknStr := r.Header.Get("Authorization")

		if tknStr == "" {

			common.ReturnResponse(w, "failed", http.StatusUnauthorized, "request does not contain an access token", nil)
			return
		}
		tknStr2 := strings.Split(tknStr, "Bearer ")
		token, err := jwt.ParseWithClaims(tknStr2[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(common.SecretKey), nil
		})
		if err != nil {
			common.ReturnResponse(w, "failed", http.StatusBadRequest, "couldn't parse token", nil)
			return
		}

		if err != nil {
			common.ReturnResponse(w, "failed", http.StatusUnauthorized, "Unauthorized user", nil)
			return

		}

		_, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			err = errors.New("couldn't parse claims")
			log.Fatal("error occurred during parsing the token", err.Error())
			common.ReturnResponse(w, "failed", http.StatusBadRequest, "couldn't parse claims", nil)
			return

		}

		if token.Valid {
			next(w, r)

		} else {
			//fmt.Println("Couldn't handle this token:", err)
			err = errors.New("couldn't handle this token")
			fmt.Printf("Couldn't handle this token:", err.Error())

			common.ReturnResponse(w, "failed", http.StatusBadRequest, "Couldn't handle this token", nil)
			return
		}

		common.ReturnResponse(w, "failed", http.StatusUnauthorized, "Unauthorized user", nil)
		return

	})
}

func ValidateJWTV2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tknStr := r.Header.Get("Authorization")

		if tknStr == "" {

			common.ReturnResponse(w, "failed", http.StatusUnauthorized, "request does not contain an access token", nil)
			return
		}
		tknStr2 := strings.Split(tknStr, "Bearer ")
		token, err := jwt.ParseWithClaims(tknStr2[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(common.SecretKey), nil
		})
		if err != nil {
			common.ReturnResponse(w, "failed", http.StatusBadRequest, "couldn't parse token", nil)
			return
		}

		if err != nil {
			common.ReturnResponse(w, "failed", http.StatusUnauthorized, "Unauthorized user", nil)
			return

		}

		_, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			err = errors.New("couldn't parse claims")
			log.Fatal("error occurred during parsing the token", err.Error())
			common.ReturnResponse(w, "failed", http.StatusBadRequest, "couldn't parse claims", nil)
			return

		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			//fmt.Println("Couldn't handle this token:", err)
			err = errors.New("couldn't handle this token")
			fmt.Printf("Couldn't handle this token:", err.Error())

			common.ReturnResponse(w, "failed", http.StatusBadRequest, "Couldn't handle this token", nil)
			return
		}
		//common.ReturnResponse(w, "failed", http.StatusUnauthorized, "Unauthorized user", nil)
		return
	})
}
