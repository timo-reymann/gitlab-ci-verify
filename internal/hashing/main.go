package hashing

import (
	"crypto/sha256"
	"encoding/hex"
)

func CreateHashFromString(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
