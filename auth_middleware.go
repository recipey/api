package main

import (
	"log"
	"net/http"
	"strings"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth_str := r.Header.Get("Authorization")
		slice_str := strings.Split(auth_str, " ")[0:1]
		bearer, token := slice_str[0], slice_str[1]

		// TODO: add function to authenticate token
		if bearer != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Request could not be authorized.")
		}

		log.Println(token)
		next.ServeHTTP(w, r)
	})
}
