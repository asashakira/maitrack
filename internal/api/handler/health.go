package handler

import (
	"fmt"
	"net/http"

	"github.com/asashakira/maitrack/internal/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking Health")
	utils.RespondWithJSON(w, 200, struct{}{})
}

func ErrorCheck(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, 400, "Something went wrong")
}
