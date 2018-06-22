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
	// docker-compose sets up the postgres database on network host named `db` for now
	connStr :=
		fmt.Sprintf("host=db user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error

	a.DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	ping_router := PingRouter{}
	ping_router.Route(a.Router)

	users_router := UsersRouter{DB: a.DB}
	users_router.Route(a.Router)

	recipes_router := RecipesRouter{DB: a.DB}
	recipes_router.Route(a.Router)

	registration_router := RegistrationRouter{DB: a.DB}
	registration_router.Route(a.Router)

	// protected endpoint requiring jwt auth
	a.Router.HandleFunc("/private", a.PrivateResource).Methods("GET", "POST", "PUT", "PATCH", "DELETE")

	a.Router.Use(loggingMiddleware)

	http.Handle("/", a.Router)
}

// Start the app
func (a *App) Run(addr string) {
	log.Println("Starting api... Listening on port", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong!")
}

// test auth in this route
func (a *App) PrivateResource(w http.ResponseWriter, r *http.Request) {
	return
}
