package jwtgenerator

import (
	"dms-api/config"
	"dms-api/utils/cryptography/security"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateOTPToken(id int) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"user_id":   	id,
		"purpose":      "password_reset",
		"verified_otp": true,
		"exp":          time.Now().Add(time.Minute * 1).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.Config("COOKIES_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	//Create Instance of NewAPISecurity
	apiSec, _ := security.NewAPISecurity()
	//Encrypt the token which is t
	encryptedToken, _ := apiSec.Encrypt(t)
	//Convert encryptedToken to Hex string
	hexToken := apiSec.BytesToHex(encryptedToken)
	// Return token
	return hexToken, nil
}
