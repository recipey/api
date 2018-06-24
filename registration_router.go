package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type RegistrationRouter struct {
	Subrouter *mux.Router
	DB        *sql.DB
}

func (rr *RegistrationRouter) Route(router *mux.Router) {
	rr.Subrouter = router.PathPrefix("/registration").Subrouter()

	rr.Subrouter.HandleFunc("", rr.RegisterUser).Methods("POST")
}

func (rr *RegistrationRouter) RegisterUser(w http.ResponseWriter, r *http.Request) {
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
