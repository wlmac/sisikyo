package oauth

import (
	"crypto/rand"
	"math/big"
)

func genRandom(l int) ([]byte, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	result := make([]byte, l)
	for i := 0; i < l; i++ {
		j, err := rand.Int(rand.Reader, new(big.Int).SetInt64(int64(len(letters))))
		if err != nil {
			return nil, err
		}
		result[i] = letters[j.Int64()]
	}
	return result, nil
}
