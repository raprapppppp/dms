package utils

import (
	"errors"
)

// Custom Error
var (
	ErrEmailNotExist       = errors.New("email does not exist")
	ErrOTPGenerationFailed = errors.New("failed to generate OTP")
	ErrNotFound            = errors.New("not found")
	ErrNotMatch            = errors.New("not match")
	ErrAlreadyExist        = errors.New("already exist")
	ErrSendingOTP          = errors.New("error sending otp")
	ErrOTPRequestLimit     = errors.New("please wait before requesting again")
	ErrNoLatestOtp		   = errors.New("no otp found")
	ErrOTPAlreadyUsed	   = errors.New("otp Already used")
	ErrOTPExpired		   = errors.New("otp expired")
	ErrOTPUpdateFailed	   = errors.New("failed to updated OTP")
	// Add other specific errors as needed, e.g., ErrSMSServiceUnavailable
)