package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/asashakira/maitrack/internal/service"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

// Context key to store the authenticated user
type contextKey string

const UserContextKey contextKey = "authenticatedUser"

// Auth is a middleware that validates JWTs.
func (m *Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// 1️⃣ Try getting token from cookie
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			tokenString = cookie.Value
		} else {
			// 2️⃣ If no cookie, try Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// If no token found in either cookies or headers, reject the request
		if tokenString == "" {
			utils.RespondWithError(w, 401, "Unauthorized")
			return
		}

		// Get the secret key
		secretKey, secretKeyErr := service.GetSecretKey()
		if secretKeyErr != nil {
			log.Println("Error fetching secret key:", secretKeyErr)
			utils.RespondWithError(w, 500, "Internal Server Error")
			return
		}

		// Parse the token
		claims := &service.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
			}
			return secretKey, nil
		})

		// Check if token is valid
		if err != nil || !token.Valid {
			log.Println("Invalid Token:", err)
			utils.RespondWithError(w, 401, "Invalid Token")
			return
		}

		log.Printf("Authenticated user: %s, Role: %s\n", claims.Username, claims.Role)

		// Store claims in context for later use
		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		next(w, r.WithContext(ctx))
	}
}
