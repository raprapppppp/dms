package otpgenerator

import (
	"crypto/rand"
	"fmt"
	"io"
)

func GenerateOTP(length int)(string, error){
	if length <= 0 {
		return "", fmt.Errorf("OTP length must be a positive integer")
	}
	const otpChars = "0123456789"
	
	otp := make([]byte, length) 

	if _, err := io.ReadFull(rand.Reader, otp); err != nil {
		return "", fmt.Errorf("failed to read random bytes for OTP: %w", err)
	}

	for i := 0; i < length; i++ {
		otp[i] = otpChars[otp[i]%byte(len(otpChars))]
	}
	return string(otp), nil
	
}