package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asashakira/mai.gg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
)

func TestCreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create pgxmock: %v", err)
	}
	defer mock.Close()

	queries := database.New(mock)
	handler := &Handler{queries: queries}

	// Test data
	userID := uuid.New()
	mockUser := database.User{
		UserID:   userID,
		SegaID:   "testSegaID",
		Password: "testPassword",
		GameName: "testGameName",
		TagLine:  "testTagLine",
	}

	// Mock queries
	// mock.ExpectQuery(`SELECT .* FROM users WHERE sega_id = \$1`).
	// 	WithArgs(mockUser.SegaID).
	// 	WillReturnRows(pgxmock.NewRows([]string{})) // User does not exist

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(mockUser.UserID, mockUser.SegaID, mockUser.Password, mockUser.GameName, mockUser.TagLine).
		WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(mockUser.UserID))

	// Create HTTP request and response recorder
	reqBody := `{"segaID":"testSegaID","password":"testPassword","gameName":"testGameName","tagLine":"testTagLine"}`
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
	// r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.CreateUser(w, r)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all database expectations were met: %v", err)
	}
}

func TestGetUserByUUID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create pgxmock: %v", err)
	}
	defer mock.Close()

	queries := database.New(mock)
	handler := &Handler{queries: queries}

	// Test data
	userID := uuid.New()
	mockUser := database.User{
		UserID:   userID,
		SegaID:   "testSegaID",
		Password: "testPassword",
		GameName: "testGameName",
		TagLine:  "testTagLine",
	}

	// Mock query
	mock.ExpectQuery(`SELECT .* FROM users WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"user_id", "sega_id", "password", "game_name", "tag_line"}).
			AddRow(mockUser.UserID, mockUser.SegaID, mockUser.Password, mockUser.GameName, mockUser.TagLine))

	// Create HTTP request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	rec := httptest.NewRecorder()

	// Set URL parameter
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("uuid", userID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	// Call handler
	handler.GetUserByID(rec, req)

	// Assert response
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all database expectations were met: %v", err)
	}
}
