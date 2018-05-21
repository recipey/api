package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize a database connection
func (a *App) Initialize(user, password, dbname string) {
	// docker-compose sets up the postgres database on network host db for now
	connStr :=
		fmt.Sprintf("host=db user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error

	a.DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	// ping
	a.Router.HandleFunc("/ping", a.pingHandler)

	// users
	a.Router.HandleFunc("/users", a.ListUsers).Methods("GET")
	a.Router.HandleFunc("/users", a.CreateUser).Methods("POST")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.GetUser).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.UpdateUser).Methods("PUT")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.DeleteUser).Methods("DELETE")

	// recipes
	a.Router.HandleFunc("/recipes", a.ListRecipes).Methods("GET")
	a.Router.HandleFunc("/recipes", a.CreateRecipe).Methods("POST")
	a.Router.HandleFunc("/recipes/{id:[0-9]+}", a.GetRecipe).Methods("GET")
	a.Router.HandleFunc("/recipes/{id:[0-9]+}", a.UpdateRecipe).Methods("PUT")
	a.Router.HandleFunc("/recipes/{id:[0-9]+}", a.DeleteRecipe).Methods("DELETE")

	http.Handle("/", a.Router)
}

// Start the app
func (a *App) Run(addr string) {
	log.Println("Starting api... listening on port", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong!")
}

func (a *App) ListRecipes(w http.ResponseWriter, r *http.Request) {
	// limit to 10 recipes for now
	recipes, err := getRecipes(a.DB, 0, 10)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, recipes)
}

func (a *App) ListUsers(w http.ResponseWriter, r *http.Request) {
	// limit to 10 users for now, probably remove limit restriction later
	users, err := getUsers(a.DB, 0, 10)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var rec recipe
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&rec); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := rec.createRecipe(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, rec)
}

func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u user
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) GetRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	rec := recipe{ID: id}
	if err := rec.getRecipe(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Recipe not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, rec)
}

func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := user{ID: id}
	if err := u.getUser(a.DB); err != nil {
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

func (a *App) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	var rec recipe
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rec); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()
	rec.ID = id

	if err := rec.updateRecipe(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, rec)
}

func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	if err := u.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	rec := recipe{ID: id}
	if err := rec.deleteRecipe(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"data": "success"})
}

func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := user{ID: id}
	if err := u.deleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"data": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
