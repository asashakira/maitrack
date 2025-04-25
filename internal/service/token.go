package service

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"sub"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GetSecretKey() ([]byte, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("JWT_SECRET not found")
	}
	return secretKey, nil
}
