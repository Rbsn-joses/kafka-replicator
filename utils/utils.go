package utils

import (
	"os"
	"unicode/utf8"
)

func GetenvOrDefaultValue(key string, defaults ...string) string {
	result := os.Getenv(key)
	if utf8.RuneCountInString(result) > 0 {
		return result
	}
	for index := 0; index < len(defaults); index++ {
		if utf8.RuneCountInString(defaults[index]) > 0 {
			return defaults[index]
		}
	}
	return ""
}
