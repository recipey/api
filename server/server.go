package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/recipey/api/user_registration"
)

// Server serves user registration routes
type Server struct {
	UserRegistration userregistration.Service // not sure if this needs to be public
	Router           *mux.Router
}

// New return new server struct to route to services
func New(ur userregistration.Service) *Server {
	r := mux.NewRouter()
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`pong!`))
	})

	// aggregating services...
	// 1. get subrouter for path
	// 2. pkg server has similar context name i.e. user_registration
	// 3. user_registration pkg defines public methods on a Service
	// 4. set routes for the subrouter to use defined services
	userRegistrationSubrouter := r.PathPrefix("/user_registration").Subrouter()
	setUserRegistrationRoutes(userRegistrationSubrouter)

	s := &Server{
		UserRegistration: ur,
		Router:           r,
	}

	return s
}
