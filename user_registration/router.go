package userregistration

import (
	"net/http"

	"github.com/gorilla/mux"
)

// TODO: not the most clean either now it's a hard dependency on gorilla/mux

// should just provide http handlers

// SetRoutes attaches routes to the given router
func SetRoutes(r mux.Router) {
	r.HandleFunc("/ping", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`pong!`))
	})
}
