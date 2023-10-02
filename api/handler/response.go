package handler

import (
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int) {
	if code > 499 {
		log.Println("Responding with 5xx error: ")
	}

	respondWithJSON(w, code)
}

func respondWithJSON(w http.ResponseWriter, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
}
