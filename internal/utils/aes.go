package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Encrypt(text string) (string, error) {
	plainText := []byte(text)

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("SECRET_KEY not found")
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plainText, nil)
	return encode(ciphertext), nil
}

func decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Decrypt(text string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("SECRET_KEY not found")
	}

	ciphertext, err := decode(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
