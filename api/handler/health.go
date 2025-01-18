package handler

import (
	"fmt"
	"net/http"

	"github.com/asashakira/mai.gg/api"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking Health")
	api.RespondWithJSON(w, 200, struct{}{})
}

func ErrorCheck(w http.ResponseWriter, r *http.Request) {
	api.RespondWithError(w, 400, "Something went wrong")
}
