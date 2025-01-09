package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/asashakira/mai.gg/api"
	"github.com/asashakira/mai.gg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	UserID    uuid.UUID        `json:"userID"`
	SegaID    string           `json:"segaID"`
	Password  string           `json:"password"`
	GameName  string           `json:"gameName"`
	TagLine   string           `json:"tagLine"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		SegaID   string `json:"segaID"`
		Password string `json:"password"`
		GameName string `json:"gameName"`
		TagLine  string `json:"tagLine"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := h.queries.CreateUser(r.Context(), database.CreateUserParams{
		UserID:   uuid.New(),
		SegaID:   params.SegaID,
		Password: params.Password,
		GameName: params.GameName,
		TagLine:  params.TagLine,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("CreateUser %v", err)
		log.Println(errorMessage)
		api.RespondWithError(w, 400, errorMessage)
		return
	}
	log.Println("CreateUser:", ConvertUser(user))
	api.RespondWithJSON(w, 200, ConvertUser(user))
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse(chi.URLParam(r, "id"))
	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		errorMessage := fmt.Sprintf("GetUserByID %v", err)
		log.Println(errorMessage)
		api.RespondWithError(w, 400, errorMessage)
		return
	}
	log.Println("GetUserByID:", ConvertUser(user))
	api.RespondWithJSON(w, 200, ConvertUser(user))
}

func (h *Handler) GetUserByMaiID(w http.ResponseWriter, r *http.Request) {
	gameName := chi.URLParam(r, "gameName")
	gameName, err := url.QueryUnescape(gameName)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("error decoding gameName from url: %v", err))
		return
	}
	tagLine := chi.URLParam(r, "tagLine")
	tagLine, err = url.QueryUnescape(tagLine)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("error decoding tagLine from url: %v", err))
		return
	}

	user, err := h.queries.GetUserByMaiID(r.Context(), database.GetUserByMaiIDParams{
		GameName: gameName,
		TagLine:  tagLine,
	})
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No user found with provided fields: %s", err)
			log.Println(errorMessage)
			api.RespondWithError(w, 404, errorMessage)
			return
		}
		errorMessage := fmt.Sprintf("GetUserByMaiID %v", err)
		log.Println(errorMessage)
		api.RespondWithError(w, 400, errorMessage)
		return
	}
	log.Println("GetUserByID:", ConvertUser(user))
	api.RespondWithJSON(w, 200, ConvertUser(user))
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.queries.GetAllUsers(r.Context())
	if err != nil {
		errorMessage := fmt.Sprintf("GetAllUsers %v", err)
		log.Println(errorMessage)
		api.RespondWithError(w, 400, errorMessage)
		return
	}
	log.Println("GetAllUsers: user count -", len(users))
	api.RespondWithJSON(w, 200, ConvertUsers(users))
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID   uuid.UUID `json:"userID,omitempty"`
		SegaID   string    `json:"segaID,omitempty"`
		Password string    `json:"password,omitempty"`
		GameName string    `json:"gameName,omitempty"`
		TagLine  string    `json:"tagLine,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	// Fetch existing user
	user, err := h.queries.GetUserByID(r.Context(), params.UserID)
	if err != nil {
		errorMessage := fmt.Sprintf("user not found %v", err)
		log.Println(errorMessage)
		api.RespondWithError(w, 400, errorMessage)
		return
	}

	// Update only the fields provided in the request
	updatedUser, err := h.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		UserID:   user.UserID,
		SegaID:   ifNotEmpty(params.SegaID, user.SegaID),
		Password: ifNotEmpty(params.Password, user.Password),
		GameName: ifNotEmpty(params.GameName, user.GameName),
		TagLine:  ifNotEmpty(params.TagLine, user.TagLine),
	})
	if err != nil {
		errorMessage := fmt.Sprintf("UpdateUser %v", err)
		log.Println(errorMessage)
		api.RespondWithError(w, 400, errorMessage)
		return
	}

	// log and reqpond updated user
	log.Println("UpdateUser:", ConvertUser(updatedUser))
	api.RespondWithJSON(w, 200, ConvertUser(updatedUser))
}

func ConvertUsers(dbUsers []database.User) []User {
	users := []User{}
	for _, user := range dbUsers {
		users = append(users, ConvertUser(user))
	}
	return users
}

func ConvertUser(dbUser database.User) User {
	return User{
		UserID:    dbUser.UserID,
		SegaID:    dbUser.SegaID,
		Password:  dbUser.Password,
		GameName:  dbUser.GameName,
		TagLine:   dbUser.TagLine,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}
