package services

import (
	"dms-api/internal/models"
	"dms-api/internal/repository"
	"dms-api/utils/cryptography/encrypt"
	"dms-api/utils/customerror"
	"dms-api/utils/otpgenerator"
	"fmt"
	"time"
)

type AuthServices interface {
	LoginService(cred models.Login) (models.User, error)
	//RegisterService(cred models.Accounts) (models.Accounts, error)
	ForgotPasswordRequestService(staffId models.Forgot) (string, error)
	VerifyOTPService(otp models.VerifyOTP) (int, error)
	UpdatePasswordService(id int, newPassword string) error
}

type InjectAuthRepository struct {
	repo repository.AuthRepository
}

func AuthServicesInit(repo repository.AuthRepository) AuthServices {
	return &InjectAuthRepository{repo}
}

// Login
func (s *InjectAuthRepository) LoginService(cred models.Login) (models.User, error) {
	isUsernameExist := s.repo.CheckUsernameIfExist(cred.Username)
	if !isUsernameExist {
		return models.User{}, customerror.ErrNotFound
	}
	//Get user single rows
	user, _ := s.repo.LoginRepo(cred.Username)
	//Hash the password before comparing
	password := encrypt.SHASecure(cred.Password)
	fmt.Print(password)
	//Check if Password is correct
	if user.UserPass != password{
		return models.User{}, customerror.ErrNotMatch
	}
	return user, nil
}

// Register
/* func (s *InjectAuthRepository) RegisterService(cred models.Accounts) (models.Accounts, error) {

	isUsernameExist := s.repo.CheckUsernameIfExist(cred.Username)
	isEmailExist := s.repo.CheckEmailIfExist(cred.Email)

	if isEmailExist || isUsernameExist {
		return models.Accounts{}, customerror.ErrAlreadyExist
	}
	cred.Password = encrypt.HashPassword(cred.Password)
	regAcc, _ := s.repo.RegisterRepo(cred)
	return regAcc, nil
} */

// ForgotPasswordRequest
func (s *InjectAuthRepository) ForgotPasswordRequestService(staffId models.Forgot) (string, error) {
	var passwordreset models.OTP
	//Check Staff ID if exist
	isStaffIdExist := s.repo.CheckStaffIdIfExist(staffId.StaffID)

	if !isStaffIdExist {
		return "", customerror.ErrStaffIDNotExist
	}
	//Get User ID via Staff ID
	userID, _ := s.repo.GetUserIDByStaffID(staffId.StaffID)

	//Rate limit
	latestOTPReq, _ := s.repo.GetLatestOTP(int(userID.UserID))
	if time.Since(latestOTPReq.CreatedAt) < 1*time.Minute {
		return "", customerror.ErrOTPRequestLimit
	}
	//Generate OTP
	otp, err := otpgenerator.GenerateOTP(6)
	if err != nil {
		return "", customerror.ErrOTPGenerationFailed
	}
	//Set accountID, OTPCode, Expiration in before saving to DB
	passwordreset.UserID = userID.UserID
	passwordreset.OTPCode = encrypt.SHASecure(otp)
	passwordreset.ExpiresAt = time.Now().Add(1 * time.Minute)
	s.repo.ForgotPasswordRequestRepo(passwordreset)

	//Send OTP in gmail
	/* mailErr := gomail.SendOTPViaMail(email.Email, otp)
	if mailErr != nil {
		return "", ErrSendingOTP
	} */
	return otp, nil
}
// VerifyOTP
func (s *InjectAuthRepository) VerifyOTPService(otp models.VerifyOTP) (int, error) {
	//Get the ID of the requestor
	userID, err := s.repo.GetUserIDByStaffID(otp.StaffID)

	if err != nil {
		return 0, customerror.ErrStaffIDNotExist
	}
	//Get the latest OTP request
	latestOTPReq, err := s.repo.GetLatestOTP(int(userID.UserID))
	if err != nil {
		return 0, customerror.ErrNoLatestOtp
	}
	//Check if the OTP is Already Used
	if latestOTPReq.Used {
		return 0, customerror.ErrOTPAlreadyUsed
	}
	//Check if OTP is expired
	if time.Now().After(latestOTPReq.ExpiresAt) {
		s.repo.UpdateOTPStatus(int(latestOTPReq.ID), true)
		return 0, customerror.ErrOTPExpired
	}
	//Hash the OTP before comparing
	hashOTP := encrypt.SHASecure(otp.Otp)
	//Check if OTP is match in DB
	if latestOTPReq.OTPCode != hashOTP {
		return 0, customerror.ErrNotMatch
	}
	//Update OTP status : Used=True
	if err := s.repo.UpdateOTPStatus(int(latestOTPReq.ID), true); err != nil {
		return 0, customerror.ErrOTPUpdateFailed
	}
	return int(userID.UserID), nil
}
//Reset the Password
func(s *InjectAuthRepository) UpdatePasswordService(id int, newPassword string) error{

	hashPassword := encrypt.SHASecure(newPassword)

	return s.repo.UpdatePasswordRepo(id, hashPassword)
}
