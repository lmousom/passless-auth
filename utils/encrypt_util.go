package utils

import (
	"crypto/sha1"
	"encoding/base64"
)

func Encrypt(message []byte) string {
	hasher := sha1.New()
	hasher.Write(message)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return sha
}
