package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"dms-api/config"
	"encoding/base64"
	"fmt"
)

func Decrypt(encodedData string) (string, error){
	key := []byte(config.Config("AES_KEY"))
	cipherData, err := base64.StdEncoding.DecodeString(encodedData)
    if err != nil {
        return "", err
    }
	// generate a new aes cipher using our 32 byte long key
	cipherBlock, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
	//Used GCM Mode for symmetric key cryptographic block ciphers
	gcm, err := cipher.NewGCM(cipherBlock)
     if err != nil {
        return "", err
    }

	nonceSize := gcm.NonceSize()
	//Ensure the input is long enough to contain nonce
    if len(cipherData) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

	//Extract nonce and ciphertext
	nonce, cipherData := cipherData[:nonceSize], cipherData[nonceSize:]

	decryptedData, err := gcm.Open(nil,nonce,cipherData,nil)
 	if err != nil {
        return "", err
    }
	fmt.Println(string(decryptedData))
	return string(decryptedData), nil
}
