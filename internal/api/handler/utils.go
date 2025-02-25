package handler

import (
	"fmt"
	"strings"
)

func ifNotNil[T any](v *T, fallback T) T {
	if v == nil {
		return fallback
	}
	return *v
}

func decodeMaiID(maiID string) (gameName, tagLine string, err error) {
	parts := strings.Split(maiID, "-")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid maiID format")
	}
	tagLine = parts[len(parts)-1]
	gameName = strings.Join(parts[:len(parts)-1], "-")
	return gameName, tagLine, nil
}
