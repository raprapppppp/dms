package encrypt

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {

	hashNaPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashNaPassword)
}

func CompareHashAndPassword(hashed string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
