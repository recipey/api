package api

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, err error) {
	respondWithJSON(w, code, map[string]error{"error": err})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
