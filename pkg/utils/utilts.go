package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashingPassword(content, salt_key string) string {
	data := salt_key + content
	hash := sha256.New()
	hash.Write([]byte(data))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func VerifyPassword(storedHash, content, salt_key string) bool {
	return HashingPassword(content, salt_key) == storedHash
}
