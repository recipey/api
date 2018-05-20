package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type UsersRouter struct {
	Subrouter *mux.Router
	DB        *sql.DB
}

func (ur *UsersRouter) Route(router *mux.Router) {
	ur.Subrouter = router.PathPrefix("/users").Subrouter()

	ur.Subrouter.HandleFunc("", ur.ListUsers).Methods("GET")
	ur.Subrouter.HandleFunc("", ur.CreateUser).Methods("POST")
	ur.Subrouter.HandleFunc("/{id:[0-9]+}", ur.GetUser).Methods("GET")
	ur.Subrouter.HandleFunc("/{id:[0-9]+}", ur.UpdateUser).Methods("PUT")
	ur.Subrouter.HandleFunc("/{id:[0-9]+}", ur.DeleteUser).Methods("DELETE")
}

func (ur *UsersRouter) ListUsers(w http.ResponseWriter, r *http.Request) {
	// limit to 10 users for now, probably remove limit restriction later
	users, err := getUsers(ur.DB, 0, 10)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (ur *UsersRouter) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := user{ID: id}
	if err := u.deleteUser(ur.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"data": "success"})
}

func (ur *UsersRouter) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u user
	decoder := json.NewDecoder(r.Body)
	// maps values from json that match back to struct
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateUser(ur.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (ur *UsersRouter) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := user{ID: id}
	if err := u.getUser(ur.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (ur *UsersRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u user
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.createUser(ur.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}
