package service

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName"`
	jwt.RegisteredClaims
}

func GetSecretKey() ([]byte, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("JWT_SECRET not found")
	}
	return secretKey, nil
}
