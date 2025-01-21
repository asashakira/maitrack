package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ifNotNil[T any](v *T, fallback T) T {
	if v == nil {
		return fallback
	}
	return *v
}

func DecodeMaiID(maiID string) (gameName, tagLine string, err error) {
	parts := strings.Split(maiID, "-")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid maiID format")
	}
	tagLine = parts[len(parts)-1]
	gameName = strings.Join(parts[:len(parts)-1], "-")
	return gameName, tagLine, nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
