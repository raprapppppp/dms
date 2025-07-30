package services

import (
	"dms-api/cryptography/encrypt"
	"dms-api/modals"
	"dms-api/otpgenerator"
	"dms-api/repository"
	"dms-api/utils"
	"time"
)


type LoginServices interface {
	LoginService(cred modals.Login) (modals.Accounts, error)
	RegisterService(cred modals.Accounts) (modals.Accounts, error)
	ForgotPasswordRequestService(email modals.Forgot) (string, error)
	VerifyOTPService(otp modals.VerifyOTP) (string, error)
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
		return modals.Accounts{}, utils.ErrNotFound
	}
	account, _ := s.repo.LoginRepo(cred.Username)
	isMatch := encrypt.CompareHashAndPassword(account.Password, cred.Password)
	if !isMatch {
		return modals.Accounts{}, utils.ErrNotMatch
	}
	return account, nil
}

// Register
func (s *InjectLoginRepository) RegisterService(cred modals.Accounts) (modals.Accounts, error) {

	isUsernameExist := s.repo.CheckUsernameIfExist(cred.Username)
	isEmailExist := s.repo.CheckEmailIfExist(cred.Email)

	if isEmailExist || isUsernameExist {
		return modals.Accounts{}, utils.ErrAlreadyExist 
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
		return "", utils.ErrEmailNotExist
	}
	//Get ID via email
	accountID, _ := s.repo.GetAccountByEmail(email.Email)

	//Rate limit
	latestOTPReq, _ := s.repo.GetLatestOTP(int(accountID.ID))
	if time.Since(latestOTPReq.CreatedAt) < 1*time.Minute {
		return "", utils.ErrOTPRequestLimit
	}
	//Generate OTP
	otp, err := otpgenerator.GenerateOTP(6)
	if err != nil {
		return "", utils.ErrOTPGenerationFailed
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

func (s *InjectLoginRepository) VerifyOTPService(otp modals.VerifyOTP) (string, error) {
	//Get the ID of the requestor
	accountID ,err  := s.repo.GetAccountByEmail(otp.Identifier)
	if err != nil {
		return "", utils.ErrNotFound
	}
	//Get the latest OTP request
	latestOTPReq,err := s.repo.GetLatestOTP(int(accountID.ID))
	if err != nil {
		return "", utils.ErrNoLatestOtp
	}
	//Check if the OTP is Already Used
	if latestOTPReq.Used {
		return "", utils.ErrOTPAlreadyUsed
	}
	if time.Now().After(latestOTPReq.ExpiresAt){
		s.repo.UpdateOTPStatus(int(latestOTPReq.ID), true)
		return "", utils.ErrOTPExpired
	}
	//Check if OTP is match in DB
	if !encrypt.CompareHashAndPassword(latestOTPReq.OTPCode, otp.Otp){
		return "", utils.ErrNotMatch
	}
	//Update OTP status : Used=True
	if err := s.repo.UpdateOTPStatus(int(latestOTPReq.ID), true); err != nil {
		return "", utils.ErrOTPUpdateFailed
	}
	return "Verified", nil

}
