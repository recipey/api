package main

import (
	"log"
	"net/http"
	"strings"
)

// steps for auth?
// 1. this is jwt auth so first read... Authorization header?
//

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth_str := r.Header.Get("Authorization")
		slice_str := strings.Split(auth_str, " ")[0:1]
		bearer, token := slice_str[0], slice_str[1]

		if bearer != "Bearer" || !authenticate_jwt(token) {
			respondWithError(w, http.StatusUnauthorized, "Request could not be authorized.")
		}

		log.Println(token)
		next.ServeHTTP(w, r)
	})
}
