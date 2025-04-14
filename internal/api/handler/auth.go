package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/asashakira/maitrack/internal/api/middleware"
	database "github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/asashakira/maitrack/internal/scraper"
	"github.com/asashakira/maitrack/internal/service"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		SegaID       string `json:"segaID"`
		SegaPassword string `json:"segaPassword"`
		GameName     string `json:"gameName"`
		TagLine      string `json:"tagLine"`
	}

	// Parse request
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, 400, "Invalid JSON payload")
		return
	}

	// Input validation
	if params.Username == "" || params.Password == "" || params.SegaID == "" || params.SegaPassword == "" {
		utils.RespondWithError(w, 400, "All fields are required")
		return
	}

	// Verify SegaID and Password
	scrapedUserData, scrapeErr := scraper.ScrapeUserData(params.SegaID, params.SegaPassword)
	if scrapeErr != nil {
		log.Printf("Failed to scrape user data: %v", scrapeErr)
		utils.RespondWithError(w, 400, "Invalid SegaID or password")
		return
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}

	// Encrypt Sega Password
	encryptedSegaPassword, err := utils.Encrypt(params.SegaPassword)
	if err != nil {
		log.Printf("Error encrypting Sega password: %v", err)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}

	// create user
	user, err := h.queries.CreateUser(r.Context(), database.CreateUserParams{
		UserID:       uuid.New(),
		Username:     params.Username,
		Password:     string(hashedPassword),
		SegaID:       params.SegaID,
		SegaPassword: encryptedSegaPassword,
		GameName:     params.GameName,
		TagLine:      params.TagLine,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Error creating user: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// create user data
	_, err = h.queries.CreateUserData(r.Context(), database.CreateUserDataParams{
		ID:              uuid.New(),
		UserID:          user.UserID,
		GameName:        user.GameName,
		TagLine:         user.TagLine,
		Rating:          scrapedUserData.Rating,
		SeasonPlayCount: scrapedUserData.SeasonPlayCount,
		TotalPlayCount:  scrapedUserData.TotalPlayCount,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Error creating user data: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// create user metadata
	defaultLastPlayedAtTime, _ := utils.StringToUTCTime("2006-01-02 15:04")
	_, err = h.queries.CreateUserMetadata(r.Context(), database.CreateUserMetadataParams{
		UserID:       user.UserID,
		LastPlayedAt: pgtype.Timestamp{Time: defaultLastPlayedAtTime, Valid: true},
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Error Creating UserMetadata: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// Generate JWT Token
	claims := service.Claims{
		Username: user.Username,
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "maitrack",
			Subject:   user.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	secretKey, err := service.GetSecretKey()
	if err != nil {
		log.Printf("Failed to get secret key: %v", err)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Failed to sign token: %v", err)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}

	// Set JWT as an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		// Secure:   true,  // Ensure Secure for HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	// Response Data
	data := map[string]any{
		"userID":          user.UserID,
		"username":        user.Username,
		"gameName":        user.GameName,
		"tagLine":         user.TagLine,
		"rating":          scrapedUserData.Rating,
		"seasonPlayCount": scrapedUserData.SeasonPlayCount,
		"totalPlayCount":  scrapedUserData.TotalPlayCount,
		"lastPlayedAt":    defaultLastPlayedAtTime,
	}

	utils.RespondWithJSON(w, 201, data)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
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
	user, err := h.queries.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		log.Println("User not found:", err)
		utils.RespondWithError(w, 400, "Invalid login credentials")
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		log.Println("Invalid password:", err)
		utils.RespondWithError(w, 401, "Invalid login credentials")
		return
	}

	// Define JWT claims
	claims := service.Claims{
		Username: user.Username,
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "maitrack",
			Subject:   user.Username,
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
		HttpOnly: true,
		// Secure:   true, // Requires HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   86400, // 1 day
	})

	// Respond with user details (excluding token)
	utils.RespondWithJSON(w, 200, map[string]any{
		"user": map[string]any{
			"username": user.Username,
			"role":     "user",
		},
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie by setting it to an expired value
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expire immediately
		MaxAge:   -1,              // Force immediate expiration
		Path:     "/",
		HttpOnly: true,
		// Secure:   true, // Requires HTTPS
		SameSite: http.SameSiteLaxMode,
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
		"username": user.Username,
		"role":     user.Role,
	}

	utils.RespondWithJSON(w, 200, data)
}
