package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"
)

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Encrypt(text string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("secretKey not found")
	}
	iv := os.Getenv("IV")
	if iv == "" {
		return "", fmt.Errorf("initialization vector not found")
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, []byte(iv))
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return Encode(cipherText), nil
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Decrypt(text string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("secretKey not found")
	}
	iv := os.Getenv("IV")
	if iv == "" {
		return "", fmt.Errorf("initialization vector not found")
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, []byte(iv))
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
