package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateLogin() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

func GeneratePassword() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func CheckPassword(password string, hash string) bool {
	return HashPassword(password) == hash
}
