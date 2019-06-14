package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// TODO: can separate into several pkgs
// main.go can then just bundle these pkgs together to form an API
type app struct {
	router *mux.Router
	db     *sql.DB
}

// Initialize a database connection
func (a *app) Initialize(user, password, dbname string) {
	// docker-compose sets up the postgres database on network host named `db` for now
	connStr :=
		fmt.Sprintf("host=db user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error

	a.db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	a.router = mux.NewRouter()

	pingRouter := pingRouter{}
	pingRouter.route(a.router)

	usersRouter := usersRouter{db: a.db}
	usersRouter.route(a.router)

	recipesRouter := recipesRouter{db: a.db}
	recipesRouter.route(a.router)

	registrationRouter := registrationRouter{db: a.db}
	registrationRouter.route(a.router)

	// protected endpoint requiring jwt auth
	a.router.HandleFunc("/private", a.privateResource).Methods("GET", "POST",
		"PUT", "PATCH", "DELETE")

	a.router.Use(loggingMiddleware)

	http.Handle("/", a.router)
}

// Start the app
func (a *app) Run(addr string) {
	log.Println("Starting api... Listening on port", addr)
	log.Fatal(http.ListenAndServe(addr, a.router))
}

func (a *app) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong!")
}

// test auth in this route
func (a *app) privateResource(w http.ResponseWriter, r *http.Request) {
	return
}
