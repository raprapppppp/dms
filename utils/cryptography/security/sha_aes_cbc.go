package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
)

func SHASecure(toHash string) string {
	hasher := sha512.New()
	hasher.Write([]byte(toHash))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

type APISecurity struct {
	iv        []byte
	key       []byte
	block     cipher.Block
	//blockMode cipher.BlockMode
}

func NewAPISecurity() (*APISecurity, error) {
	key := []byte("0123456789qwerty") // 16 bytes (AES-128)
	iv := []byte("qwerty9876543210")  // 16 bytes IV (should match key length)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &APISecurity{
		iv:    iv,
		key:   key,
		block: block,
	}, nil
}

func (a *APISecurity) padString(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{0}, padding) // Java-style padding (zeroes)
	return append(src, padtext...)
}

func (a *APISecurity) unpadString(src []byte) []byte {
	// remove trailing zeroes
	return bytes.TrimRight(src, "\x00")
}

func (a *APISecurity) Encrypt(text string) ([]byte, error) {
	if len(text) == 0 {
		return nil, errors.New("empty string")
	}

	plaintext := a.padString([]byte(text))
	ciphertext := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(a.block, a.iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

func (a *APISecurity) Decrypt(hexString string) ([]byte, error) {
	if len(hexString) == 0 {
		return nil, errors.New("empty string")
	}

	data, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	if len(data)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	decrypted := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(a.block, a.iv)
	mode.CryptBlocks(decrypted, data)

	return a.unpadString(decrypted), nil
}

func (a *APISecurity) BytesToHex(data []byte) string {
	return hex.EncodeToString(data)
}
func (a *APISecurity) HexToBytes(data string) ([]byte, error) {
    decoded, err := hex.DecodeString(data)
    if err != nil {
        return nil, fmt.Errorf("error decoding hex: %w", err)
    }
    return decoded, nil
}