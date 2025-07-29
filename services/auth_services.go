package services

import (
	"dms-api/cryptography/encrypt"
	"dms-api/gomail"
	"dms-api/modals"
	"dms-api/otpgenerator"
	"dms-api/repository"
	"errors"
	"time"
)

type LoginServices interface {
	LoginService(cred modals.Login) (modals.Accounts, error)
	RegisterService(cred modals.Accounts) (modals.Accounts, error)
	ForgotPasswordRequestService(email modals.Forgot) (string, error)
}

type InjectLoginRepository struct {
	repo repository.LoginRepository
}

func LoginServicesInit(repo repository.LoginRepository) LoginServices {
	return &InjectLoginRepository{repo}
}

// Custom Error
var (
	ErrEmailNotExist       = errors.New("email does not exist")
	ErrOTPGenerationFailed = errors.New("failed to generate OTP")
	ErrNotFound            = errors.New("not found")
	ErrNotMatch            = errors.New("not match")
	ErrAlreadyExist        = errors.New("already exist")
	ErrSendingOTP          = errors.New("error sending otp")
	ErrOTPRequestLimit     = errors.New("please wait before requesting again")
	// Add other specific errors as needed, e.g., ErrSMSServiceUnavailable
)

// Login
func (s *InjectLoginRepository) LoginService(cred modals.Login) (modals.Accounts, error) {
	isUsernameExist := s.repo.CheckUsernameIfExist(cred.Username)
	if !isUsernameExist {
		return modals.Accounts{}, ErrNotFound
	}
	account, _ := s.repo.LoginRepo(cred.Username)
	isMatch := encrypt.CompareHashAndPassword(account.Password, cred.Password)
	if !isMatch {
		return modals.Accounts{}, ErrNotMatch
	}
	return account, nil
}

// Register
func (s *InjectLoginRepository) RegisterService(cred modals.Accounts) (modals.Accounts, error) {

	isUsernameExist := s.repo.CheckUsernameIfExist(cred.Username)
	isEmailExist := s.repo.CheckEmailIfExist(cred.Email)

	if isEmailExist || isUsernameExist {
		return modals.Accounts{}, ErrAlreadyExist
	}
	cred.Password = encrypt.HashPassword(cred.Password)
	regAcc, _ := s.repo.RegisterRepo(cred)
	return regAcc, nil
}

// ForgotPasswordRequest
func (s *InjectLoginRepository) ForgotPasswordRequestService(email modals.Forgot) (string, error) {
	var passwordreset modals.PasswordReset
	//Check email if exist
	isEmailExist := s.repo.CheckEmailIfExist(email.Email)
	if !isEmailExist {
		return "", ErrEmailNotExist
	}
	//Get ID via email
	account, _ := s.repo.GetAccountByEmail(email.Email)

	//Rate limit
	latestRequest, _ := s.repo.GetLatestCreateTime(int(account.ID))
	if time.Since(latestRequest.CreatedAt) < 1*time.Minute {
		return "", ErrOTPRequestLimit
	}
	//Generate OTP
	otp, err := otpgenerator.GenerateOTP(6)
	if err != nil {
		return "", ErrOTPGenerationFailed
	}
	//Set accountID, OTPCode, Expiration in before saving to DB
	passwordreset.AccountsID = account.ID
	passwordreset.OTPCode = encrypt.HashPassword(otp)
	passwordreset.ExpiresAt = time.Now().Add(1 * time.Minute)
	s.repo.ForgotPasswordRequestRepo(passwordreset)

	//Send OTP in gmail
	mailErr := gomail.SendOTPViaMail(email.Email, otp)
	if mailErr != nil {
		return "", ErrSendingOTP
	}
	return "sent", nil
}
