package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type recipesRouter struct {
	subrouter *mux.Router
	db        *sql.DB
}

func (rr *recipesRouter) route(router *mux.Router) {
	rr.subrouter = router.PathPrefix("/recipes").Subrouter()

	rr.subrouter.HandleFunc("", rr.listRecipes).Methods("GET")
	rr.subrouter.HandleFunc("", rr.createRecipe).Methods("POST")
	rr.subrouter.HandleFunc("/{id:[0-9]+}", rr.getRecipe).Methods("GET")
	rr.subrouter.HandleFunc("/{id:[0-9]+}", rr.updateRecipe).Methods("PUT")
	rr.subrouter.HandleFunc("/{id:[0-9]+}", rr.deleteRecipe).Methods("DELETE")

	// HARDCODED
	rr.Subrouter.HandleFunc("/list", rr.listHardCoded).Methods("GET")
	rr.Subrouter.HandleFunc("/single", rr.getHardCoded).Methods("GET")
}

func (rr *recipesRouter) listRecipes(w http.ResponseWriter, r *http.Request) {
	// limit to 10 recipes for now
	recipes, err := getRecipes(rr.DB, 0, 10)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
	}

	respondWithJSON(w, http.StatusOK, recipes)
}

func (rr *recipesRouter) createRecipe(w http.ResponseWriter, r *http.Request) {
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

func (rr *recipesRouter) getRecipe(w http.ResponseWriter, r *http.Request) {
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

func (rr *recipesRouter) updateRecipe(w http.ResponseWriter, r *http.Request) {
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

func (rr *recipesRouter) deleteRecipe(w http.ResponseWriter, r *http.Request) {
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

func (rr *recipesRouter) listHardCoded(w http.ResponseWriter, r *http.Request) {
	data := []byte(`{
		"recipes": [{
			"recipeId": 1,
			"recipeImageUrl": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSCQaBD-kga2C-NKq0WSRWAPouX-dNFeCKntpHrKhs58hV36t93",
			"recipeName": "Seared Ahi Tuna Tacos With Guacamole"
		}, {
			"recipeId": 2,
			"recipeImageUrl": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSCQaBD-kga2C-NKq0WSRWAPouX-dNFeCKntpHrKhs58hV36t93",
			"recipeName": "Chicken Tikki Masala"
		}]
	}
	`)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (rr *recipesRouter) getHardCoded(w http.ResponseWriter, r *http.Request) {
	data := []byte(`
		{
			"recipeId": 1,
			"recipeImageUrl": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSCQaBD-kga2C-NKq0WSRWAPouX-dNFeCKntpHrKhs58hV36t93",
			"recipeName": "Seared Ahi Tuna Tacos With Guacamole"
		}
	`)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
