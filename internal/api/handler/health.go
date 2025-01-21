package handler

import (
	"fmt"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking Health")
	respondWithJSON(w, 200, struct{}{})
}

func ErrorCheck(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}
