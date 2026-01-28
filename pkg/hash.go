package pkg

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

func HashString(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func CompareHash(hash string, value string) bool {
	hashed := HashString(value)
	return subtle.ConstantTimeCompare([]byte(hash), []byte(hashed)) == 1
}
