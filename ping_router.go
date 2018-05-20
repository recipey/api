package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type PingRouter struct {
	Subrouter *mux.Router
}

func (p *PingRouter) Route(router *mux.Router) {
	p.Subrouter = router.PathPrefix("/ping").Subrouter().StrictSlash(true)

	p.Subrouter.HandleFunc("/", p.PingHandler).Methods("GET")
}

func (p *PingRouter) PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong!")
}
