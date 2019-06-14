package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type registrationRouter struct {
	subrouter *mux.Router
	db        *sql.DB
}

func (rr *registrationRouter) route(router *mux.Router) {
	rr.subrouter = router.PathPrefix("/registration").Subrouter()

	rr.subrouter.HandleFunc("", rr.RegisterUser).Methods("POST")
}

func (rr *registrationRouter) registerUser(w http.ResponseWriter, r *http.Request) {
	var u user

	// create decoder that reads from request body stream
	decoder := json.NewDecoder(r.Body)

	// calling decode onto user struct
	// NOTE: Decode will only work if the user struct maps the json key to the struct field
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid request payload"))
		return
	}

	if err := u.createUser(rr.DB); err != nil {
		// hmm need to respond with either StatusBadRequest or StatusInternalServerError
		// how does the error object work and can it be distinguished?
	}

	// TODO: this should respond with json containing jwt
	respondWithJSON(w, http.StatusCreated, u)
}
