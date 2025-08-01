package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"dms-api/config"
	"encoding/base64"
	"fmt"
	"io"
)

func Encrypt(toEncrypt string) (string) {

	text := []byte(toEncrypt)
	key := []byte(config.Config("AES_KEY"))
	// generate a new aes cipher using our 32 byte long key
    cipherBlock, err := aes.NewCipher(key)
    if err != nil {
        fmt.Println(err)
    }
	//Used GCM Mode for symmetric key cryptographic block ciphers
	gcm, err := cipher.NewGCM(cipherBlock)
    if err != nil {
        fmt.Println(err)
    }
	// creates a new byte array the size of the nonce|IV
    // which must be passed to Seal
    nonce := make([]byte, gcm.NonceSize())
	// populates our nonce|IV with a cryptographically secure
    // random sequence
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        fmt.Println(err)
    }
	//Instead storing nonce|IV to a different storage | put nonce in dst arguments
	cipherData := gcm.Seal(nonce,nonce,text,nil)
	// Encode ciphertext as Base64 string
    encodedData := base64.StdEncoding.EncodeToString(cipherData)

	return encodedData
}