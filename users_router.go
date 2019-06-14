package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// TODO: review DDD
// Create a datastore pkg
// Routers should not know how to iface with sql or w/e storage
// One step further is that datastore should not know logic about "domain"
// User domain, Recipe domain, etc
// If any of these domains need add'l storage or caching then it should be handled there
type usersRouter struct {
	subrouter *mux.Router
	db        *sql.DB
}

func (ur *usersRouter) route(router *mux.Router) {
	ur.subrouter = router.PathPrefix("/users").Subrouter()

	ur.subrouter.HandleFunc("", ur.listUsers).Methods("GET")
	ur.subrouter.HandleFunc("", ur.createUser).Methods("POST")
	ur.subrouter.HandleFunc("/{id:[0-9]+}", ur.getUser).Methods("GET")
	ur.subrouter.HandleFunc("/{id:[0-9]+}", ur.updateUser).Methods("PUT")
	ur.subrouter.HandleFunc("/{id:[0-9]+}", ur.deleteUser).Methods("DELETE")
}

func (ur *usersRouter) listUsers(w http.ResponseWriter, r *http.Request) {
	// limit to 10 users for now, probably remove limit restriction later
	users, err := getUsers(ur.db, 0, 10)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (ur *usersRouter) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid user ID"))
		return
	}

	u := user{ID: id}
	if err := u.deleteUser(ur.db); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"data": "success"})
}

func (ur *usersRouter) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid user ID"))
		return
	}

	var u user
	decoder := json.NewDecoder(r.Body)
	// maps values from json that match back to struct
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid request payload"))
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateUser(ur.db); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (ur *usersRouter) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid user ID"))
		return
	}

	u := user{ID: id}
	if err := u.getUser(ur.db); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, errors.New("User not found"))
		default:
			respondWithError(w, http.StatusInternalServerError, err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (ur *usersRouter) createUser(w http.ResponseWriter, r *http.Request) {
	var u user
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid request payload"))
		return
	}
	defer r.Body.Close()

	if err := u.createUser(ur.db); err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}
