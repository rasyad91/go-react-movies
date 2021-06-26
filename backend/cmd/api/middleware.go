package main

import (
	"backend/internal/util"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allows localhost:3000(frontend) <-> localhost:4000(backend)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// allows change to header
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		next.ServeHTTP(w, r)
	})
}

func checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Before adding Vary, Authorization: %#v\n", w.Header())
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			// could set an anonymous user
		}

		fmt.Printf("After adding vary authorization: %#v\n", w.Header())

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			util.ErrorJSON(w, errors.New("invalid auth header"))
			return
		}

		if headerParts[0] != "Bearer" {
			util.ErrorJSON(w, errors.New("unauthorized - no bearer"))
			return
		}

		token := headerParts[1]
		claims, err := jwt.HMACCheck([]byte(token), []byte(cfg.JWT.Secret))
		if err != nil {
			util.ErrorJSON(w, errors.New("unauthorized - failed hmac check"))
			return
		}

		if !claims.Valid(time.Now()) {
			util.ErrorJSON(w, errors.New("unauthorized - token expired"))
			return
		}

		if !claims.AcceptAudience("mydomain.com") {
			util.ErrorJSON(w, errors.New("unauthorized - invalid audience"))
			return
		}

		if claims.Issuer != "mydomain.com" {
			util.ErrorJSON(w, errors.New("unauthorized - invalid issuer"))
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			util.ErrorJSON(w, errors.New("unauthorized"))
			return
		}
		log.Println(userID)

		next.ServeHTTP(w, r)
	})
}
