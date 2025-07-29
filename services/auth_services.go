package services

import (
	"dms-api/cryptography/encrypt"
	"dms-api/modals"
	"dms-api/otpgenerator"
	"dms-api/repository"
	"errors"
	"fmt"
)

type LoginServices interface {
	LoginService(cred modals.Login) (modals.Accounts, string)
	RegisterService(cred modals.Accounts) (modals.Accounts, string)
	ForgotPasswordService(email modals.Forgot) (string, error)
}

type InjectLoginRepository struct {
	repo repository.LoginRepository
}

func LoginServicesInit (repo repository.LoginRepository) LoginServices {
	return &InjectLoginRepository{repo}
}
//Custom Error
var (
	ErrEmailNotExist       = errors.New("email does not exist")
	ErrOTPGenerationFailed = errors.New("failed to generate OTP")
	ErrNotFound			   = errors.New("not found")
	ErrNotMatch			   = errors.New("not match")
	// Add other specific errors as needed, e.g., ErrSMSServiceUnavailable
)
//Login
func(s *InjectLoginRepository) LoginService(cred modals.Login) (modals.Accounts, string) {
	isUsernameExist := s.repo.CheckUsernameIfExist(cred.Username)
	if !isUsernameExist {
		return modals.Accounts{}, "NotFound"
	}
	account, _ := s.repo.LoginRepo(cred.Username)

	isMatch := encrypt.CompareHashAndPassword(account.Password,cred.Password)
	if !isMatch {
		return modals.Accounts{}, "NotMatch"
	}
	return account, "Match"
}
//Register
func(s *InjectLoginRepository) RegisterService(cred modals.Accounts)(modals.Accounts, string){

	isUsernameExist := s.repo.CheckUsernameIfExist(cred.Username)
	isEmailExist := s.repo.CheckEmailIfExist(cred.Email)

	if isEmailExist || isUsernameExist {
		return modals.Accounts{}, "Already Exist"
	}
	cred.Password = encrypt.HashPassword(cred.Password)
	regAcc, _ := s.repo.RegisterRepo(cred)
	return regAcc, ""
}
//ForgotPassword
func(s *InjectLoginRepository) ForgotPasswordService(email modals.Forgot)(string, error){
	isEmailExist := s.repo.CheckEmailIfExist(email.Email)
	if !isEmailExist {
		return "", ErrEmailNotExist
	}

	otp, err := otpgenerator.GenerateOTP(6)
	if err != nil {
		return "", ErrOTPGenerationFailed
	}
	fmt.Print(otp)
	return otp, nil

}

