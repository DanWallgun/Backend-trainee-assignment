package utils

import (
	"crypto/rand"
	"strings"
)

const charSet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateRandomString - generates string with given length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	sb := new(strings.Builder)
	for i := 0; i < length; i++ {
		sb.WriteByte(charSet[bytes[i]%62])
	}

	return sb.String(), nil
}
