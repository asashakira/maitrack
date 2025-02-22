package middleware

import (
	"net/http"
)

func (m *Middleware) APIKeyAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement

		// apiKey := r.Header.Get("X-API-Key")
		// if apiKey == "" {
		// 	utils.RespondWithError(w, 401, "invalid API key")
		// 	return
		// }
		next(w, r)
	}
}
