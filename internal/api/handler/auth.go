package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/asashakira/maitrack/internal/api/middleware"
	database "github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/asashakira/maitrack/internal/service"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/asashakira/maitrack/pkg/maimaiclient"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID       string `json:"userID"`
		DisplayName  string `json:"displayName"`
		Password     string `json:"password"`
		SegaID       string `json:"segaID"`
		SegaPassword string `json:"segaPassword"`
	}

	// Parse request
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, 400, "Invalid JSON payload")
		return
	}

	// Input validation
	// TODO: do better
	if params.UserID == "" || params.Password == "" || params.SegaID == "" || params.SegaPassword == "" {
		utils.RespondWithError(w, 400, "All fields are required")
		return
	}

	// Try to Login to maimaidx.net to verify
	m := maimaiclient.New()
	err := m.Login(params.SegaID, params.SegaPassword)
	if err != nil {
		log.Printf("Invalid SEGA Credentials '%s': %s\n", params.SegaID, err)
		utils.RespondWithError(w, 400, "Invalid SEGA Credentials")
		return
	}

	// Hash Password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}

	// Encrypt SEGA Creds
	encryptedSegaID, err := utils.Encrypt(params.SegaID)
	if err != nil {
		log.Printf("Error encrypting SEGA ID: %s", err)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}
	encryptedSegaPassword, err := utils.Encrypt(params.SegaPassword)
	if err != nil {
		log.Printf("Error encrypting SEGA password: %s", err)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}

	defaultTime, _ := utils.StringToUTCTime("2006-01-02 15:04")
	// create user
	user, err := h.queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:                    uuid.New(),
		UserID:                params.UserID,
		DisplayName:           params.DisplayName,
		PasswordHash:          string(passwordHash),
		EncryptedSegaID:       encryptedSegaID,
		EncryptedSegaPassword: encryptedSegaPassword,
		LastPlayedAt:          pgtype.Timestamp{Time: defaultTime, Valid: true},
		LastScrapedAt:         pgtype.Timestamp{Time: defaultTime, Valid: true},
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Error creating user: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// create initial userdata
	_, createUserDataErr := h.queries.CreateUserData(r.Context(), database.CreateUserDataParams{
		ID:              uuid.New(),
		UserUuid:        user.ID,
		Rating:          0,
		SeasonPlayCount: 0,
		TotalPlayCount:  0,
	})
	if createUserDataErr != nil {
		errorMessage := fmt.Sprintf("failed to create user data for user '%s': %s", user.UserID, createUserDataErr)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// create initial user metadata
	_, err = h.queries.CreateUserMetadata(r.Context(), database.CreateUserMetadataParams{
		UserUuid: user.ID,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Error Creating UserMetadata: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// Define JWT claims
	claims := service.Claims{
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "maitrack",
			Subject:   user.UserID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
		},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey, err := service.GetSecretKey()
	if err != nil {
		log.Println("Failed to get secret key:", err)
		utils.RespondWithError(w, 500, "Internal server error")
		return
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Failed to generate token:", err)
		utils.RespondWithError(w, 500, "Internal server error")
		return
	}

	// Set token as an HTTP-only secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	// Response Data
	data := map[string]any{
		"id":          user.ID,
		"userID":      user.UserID,
		"displayName": user.DisplayName,
	}

	utils.RespondWithJSON(w, 201, data)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID   string `json:"userID"`
		Password string `json:"password"`
	}

	// Decode request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	// Fetch user from DB
	user, err := h.queries.GetPasswordHashByUserID(r.Context(), params.UserID)
	if err != nil {
		log.Println("User not found:", err)
		utils.RespondWithError(w, 400, "Invalid login credentials")
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		log.Println("Invalid password:", err)
		utils.RespondWithError(w, 401, "Invalid login credentials")
		return
	}

	// Define JWT claims
	claims := service.Claims{
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "maitrack",
			Subject:   user.UserID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
		},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey, err := service.GetSecretKey()
	if err != nil {
		log.Println("Failed to get secret key:", err)
		utils.RespondWithError(w, 500, "Internal server error")
		return
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Failed to generate token:", err)
		utils.RespondWithError(w, 500, "Internal server error")
		return
	}

	// Set token as an HTTP-only secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	// Respond with user details (excluding token)
	utils.RespondWithJSON(w, 200, map[string]any{
		"user": map[string]any{
			"userID":      user.UserID,
			"displayName": user.DisplayName,
		},
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie by setting it to an expired value
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,              // Force immediate expiration
		Expires:  time.Unix(0, 0), // Expire immediately
	})

	utils.RespondWithJSON(w, 200, map[string]string{
		"message": "Logout successful",
	})
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*service.Claims)
	if !ok {
		fmt.Println("no authenticated user found in context")
	}

	data := map[string]any{
		"userID":      user.UserID,
		"displayName": user.DisplayName,
	}

	utils.RespondWithJSON(w, 200, data)
}
