package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func RemoveFromString(input, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, "")
}

func FormatDate(s string) string {
	return strings.ReplaceAll(s, "/", "-")
}

func StringToInt32(s string) (int32, error) {
	// Convert string to int
	value, err := strconv.Atoi(s) // Converts to int
	if err != nil {
		return 0, err
	}

	// Check if the value fits within int32 range
	if value < -2147483648 || value > 2147483647 {
		return 0, fmt.Errorf("value out of int32 range")
	}

	// Convert to int32
	return int32(value), nil
}

// create alt key using song title and artist
func CreateAltKey(title, artist string) string {
	altkey := title + artist
	altkey = strings.ToLower(altkey) // all lowercase
	return RemoveFromString(altkey, `[^一-龠ぁ-ゔァ-ヴーa-zA-Z0-9ａ-ｚＡ-Ｚ０-９々〆〤ヶ]+`)
}
