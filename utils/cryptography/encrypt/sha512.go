package encrypt

import (
	"crypto/sha512"
	"encoding/hex"
)

func SHASecure(toHash string) string {
	hasher := sha512.New()
	hasher.Write([]byte(toHash))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
