package jwtgenerator

import (
	"dms-api/config"
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
	// Return token
	return t, nil
}
