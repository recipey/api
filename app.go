package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize a database connection
func (a *App) Initialize(user, password, dbname string) {
	connStr :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error

	a.DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/ping", pingHandler)

	http.Handle("/", a.Router)
}

// Start the app
func (a *App) Run(addr string) {
	log.Println("Starting api... listening on port", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong!")
}
