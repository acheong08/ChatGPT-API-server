package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
