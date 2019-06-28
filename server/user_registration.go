package server

import (
	// "context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/recipey/api/user_registration"
)

type userRegistrationHandler struct {
	s userregistration.Service
}

func setUserRegistrationRoutes(s *mux.Router) {
	s.HandleFunc("/ping", pingHandler)
	s.HandleFunc("", registerHandler).Methods(http.MethodPost)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`pong!`))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`register!`))
}

// func (h *userRegistrationHandler) registerUser(w http.ResponseWriter, r *http.Request) {
//   ctx := context.Background()
//
//   w.Header().Set("Content-Type", "application/json; charset=utf-8")
// }
