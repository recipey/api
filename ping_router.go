package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type pingRouter struct {
	subrouter *mux.Router
}

func (p *pingRouter) route(router *mux.Router) {
	p.subrouter = router.PathPrefix("/ping").Subrouter().StrictSlash(true)

	p.subrouter.HandleFunc("/", p.pingHandler).Methods("GET")
}

func (p *pingRouter) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong!")
}
