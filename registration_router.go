package main

import (
	"database/sql"
	"encoding/json"
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
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	errors := map[string]string{}

	// check for required arguments
	if u.Username == "" {
		errors["username"] = "username is required"
	}

	if u.Email == "" {
		errors["email"] = "email is required"
	}

	if u.Password == "" {
		errors["password"] = "password is required"
	}

	if u.Password != u.PasswordConfirmation {
		errors["password_confirmation"] = "password confirmation does not match password"
	}

	if len(errors) > 0 {
		json_str, err := json.Marshal(errors)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			err_msg := string(json_str[:])
			respondWithError(w, http.StatusBadRequest, err_msg)
		}

		return
	}

	// TODO: this should respond with json containing jwt
	respondWithJSON(w, http.StatusCreated, u)
}
