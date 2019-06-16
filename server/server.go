package server

import (
	"github.com/gorilla/mux"

	"github.com/recipey/api/user_registration"
)

// Server serves user registration routes
type Server struct {
	UserRegistration userregistration.Service

	router mux.Router
}

// New TODO...
func New(ur userregistration.Service) *Server {
	s := &Server{
		UserRegistration: ur,
	}

	// r := mux.NewRouter()
	// s.UserRegistration.SetRoutes(r.PathPrefix("/user_registration").Subrouter())

	return s
}
