package services

import (
	"dms-api/internal/modals"
	"dms-api/internal/repository"
	"dms-api/utils/cryptography/encrypt"
	"dms-api/utils/customerror"
	"dms-api/utils/otpgenerator"
	"time"
)

type LoginServices interface {
	LoginService(cred modals.Login) (modals.Accounts, error)
	RegisterService(cred modals.Accounts) (modals.Accounts, error)
	ForgotPasswordRequestService(email modals.Forgot) (string, error)
	VerifyOTPService(otp modals.VerifyOTP) (int, error)
}

type InjectLoginRepository struct {
	repo repository.LoginRepository
}

func LoginServicesInit(repo repository.LoginRepository) LoginServices {
	return &InjectLoginRepository{repo}
}

// Login
func (s *InjectLoginRepository) LoginService(cred modals.Login) (modals.Accounts, error) {
	isUsernameExist := s.repo.CheckUsernameIfExist(cred.Username)
	if !isUsernameExist {
		return modals.Accounts{}, customerror.ErrNotFound
	}
	account, _ := s.repo.LoginRepo(cred.Username)
	isMatch := encrypt.CompareHashAndPassword(account.Password, cred.Password)
	if !isMatch {
		return modals.Accounts{}, customerror.ErrNotMatch
	}
	return account, nil
}

// Register
func (s *InjectLoginRepository) RegisterService(cred modals.Accounts) (modals.Accounts, error) {

	isUsernameExist := s.repo.CheckUsernameIfExist(cred.Username)
	isEmailExist := s.repo.CheckEmailIfExist(cred.Email)

	if isEmailExist || isUsernameExist {
		return modals.Accounts{}, customerror.ErrAlreadyExist
	}
	cred.Password = encrypt.HashPassword(cred.Password)
	regAcc, _ := s.repo.RegisterRepo(cred)
	return regAcc, nil
}

// ForgotPasswordRequest
func (s *InjectLoginRepository) ForgotPasswordRequestService(email modals.Forgot) (string, error) {
	var passwordreset modals.OTP
	//Check email if exist
	isEmailExist := s.repo.CheckEmailIfExist(email.Email)
	if !isEmailExist {
		return "", customerror.ErrEmailNotExist
	}
	//Get ID via email
	accountID, _ := s.repo.GetAccountByEmail(email.Email)

	//Rate limit
	latestOTPReq, _ := s.repo.GetLatestOTP(int(accountID.ID))
	if time.Since(latestOTPReq.CreatedAt) < 1*time.Minute {
		return "", customerror.ErrOTPRequestLimit
	}
	//Generate OTP
	otp, err := otpgenerator.GenerateOTP(6)
	if err != nil {
		return "", customerror.ErrOTPGenerationFailed
	}
	//Set accountID, OTPCode, Expiration in before saving to DB
	passwordreset.AccountsID = accountID.ID
	passwordreset.OTPCode = encrypt.HashPassword(otp)
	passwordreset.ExpiresAt = time.Now().Add(1 * time.Minute)
	s.repo.ForgotPasswordRequestRepo(passwordreset)

	//Send OTP in gmail
	/* mailErr := gomail.SendOTPViaMail(email.Email, otp)
	if mailErr != nil {
		return "", ErrSendingOTP
	} */
	return otp, nil
}

func (s *InjectLoginRepository) VerifyOTPService(otp modals.VerifyOTP) (int, error) {
	//Get the ID of the requestor
	accountID, err := s.repo.GetAccountByEmail(otp.Identifier)
	if err != nil {
		return 0, customerror.ErrNotFound
	}
	//Get the latest OTP request
	latestOTPReq, err := s.repo.GetLatestOTP(int(accountID.ID))
	if err != nil {
		return 0, customerror.ErrNoLatestOtp
	}
	//Check if the OTP is Already Used
	if latestOTPReq.Used {
		return 0, customerror.ErrOTPAlreadyUsed
	}
	if time.Now().After(latestOTPReq.ExpiresAt) {
		s.repo.UpdateOTPStatus(int(latestOTPReq.ID), true)
		return 0, customerror.ErrOTPExpired
	}
	//Check if OTP is match in DB
	if !encrypt.CompareHashAndPassword(latestOTPReq.OTPCode, otp.Otp) {
		return 0, customerror.ErrNotMatch
	}
	//Update OTP status : Used=True
	if err := s.repo.UpdateOTPStatus(int(latestOTPReq.ID), true); err != nil {
		return 0, customerror.ErrOTPUpdateFailed
	}
	return int(accountID.ID), nil
}
