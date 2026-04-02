package service

import (
	"crypto/sha256"
	"encoding/hex"
)

func hashSHA256Password(password string) string {
	sum := sha256.Sum256([]byte(password))
	return sha256PasswordPrefix + hex.EncodeToString(sum[:])
}
