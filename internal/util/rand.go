package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const alphanumericCharacterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString generate a random string using the given character set
func RandomString(characterSet string, length int) (string, error) {
	charSetLength := big.NewInt(int64(len(characterSet)))

	var keyBuilder strings.Builder
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charSetLength)
		if err != nil {
			return "", fmt.Errorf("generating random number: %w", err)
		}
		keyBuilder.WriteRune(rune(characterSet[randomIndex.Int64()]))
	}
	return keyBuilder.String(), nil
}

// RandomAlphanumericString generate a random string using a-z, A-Z, and 0-9
func RandomAlphanumericString(length int) (string, error) {
	return RandomString(alphanumericCharacterSet, length)
}
