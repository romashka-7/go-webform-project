package security

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}
