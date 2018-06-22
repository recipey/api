package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type RecipesRouter struct {
	Subrouter *mux.Router
	DB        *sql.DB
}

func (rr *RecipesRouter) Route(router *mux.Router) {
	rr.Subrouter = router.PathPrefix("/recipes").Subrouter()

	rr.Subrouter.HandleFunc("", rr.ListRecipes).Methods("GET")
	rr.Subrouter.HandleFunc("", rr.CreateRecipe).Methods("POST")
	rr.Subrouter.HandleFunc("/{id:[0-9]+}", rr.GetRecipe).Methods("GET")
	rr.Subrouter.HandleFunc("/{id:[0-9]+}", rr.UpdateRecipe).Methods("PUT")
	rr.Subrouter.HandleFunc("/{id:[0-9]+}", rr.DeleteRecipe).Methods("DELETE")
}

func (rr *RecipesRouter) ListRecipes(w http.ResponseWriter, r *http.Request) {
	// limit to 10 recipes for now
	recipes, err := getRecipes(rr.DB, 0, 10)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
	}

	respondWithJSON(w, http.StatusOK, recipes)
}

func (rr *RecipesRouter) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var rec recipe
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&rec); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid request payload"))
		return
	}
	defer r.Body.Close()

	if err := rec.createRecipe(rr.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, rec)
}

func (rr *RecipesRouter) GetRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid recipe ID"))
		return
	}

	rec := recipe{ID: id}
	if err := rec.getRecipe(rr.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, errors.New("Recipe not found"))
		default:
			respondWithError(w, http.StatusInternalServerError, err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, rec)
}

func (rr *RecipesRouter) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid recipe ID"))
		return
	}

	var rec recipe
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rec); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	defer r.Body.Close()
	rec.ID = id

	if err := rec.updateRecipe(rr.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, rec)
}

func (rr *RecipesRouter) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid recipe ID"))
		return
	}

	rec := recipe{ID: id}
	if err := rec.deleteRecipe(rr.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"data": "success"})
}
