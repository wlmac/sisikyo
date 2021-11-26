// Package util provides utilities.
package util

import (
	"crypto/rand"
	"math/big"
)

// GenRandom generates a base64 string using crypto/rand.
func GenRandom(l int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	result := make([]byte, l)
	for i := 0; i < l; i++ {
		j, err := rand.Int(rand.Reader, new(big.Int).SetInt64(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[j.Int64()]
	}
	return string(result), nil
}
