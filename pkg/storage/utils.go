package storage

import (
	"crypto/rand"
	"encoding/hex"
)

func randomString(len int) (string, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
