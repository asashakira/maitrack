package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/asashakira/mai.gg/internal/service"
	"github.com/asashakira/mai.gg/internal/utils"
	"github.com/golang-jwt/jwt/v5"
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

// Auth is a middleware that validates JWTs.
func (m *Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.RespondWithError(w, 401, "Unauthorized")
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		secretKey, secretKeyErr := service.GetSecretKey()
		if secretKeyErr != nil {
			errorMessage := fmt.Sprintf("%s", secretKeyErr)
			log.Println(errorMessage)
			utils.RespondWithError(w, 500, errorMessage)
			return
		}

		// Parse the token with custom claims
		claims := &service.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token was signed with HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			fmt.Println(claims)
			return secretKey, nil
		})
		if err != nil {
			log.Println("Invalid Token:", err)
			utils.RespondWithError(w, 401, "Invalid Token")
			return
		}

		if !token.Valid {
			utils.RespondWithError(w, 401, "Invalid Token")
			return
		}

		log.Printf("Authenticated user: %s, Role: %s\n", claims.Username, claims.Role)

		next(w, r)
	}
}
