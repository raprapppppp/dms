package repository

import (
	"dms-api/modals"
	"errors"

	"gorm.io/gorm"
)

type LoginRepository interface {
	LoginRepo(username string) (modals.Accounts, error)
	RegisterRepo(cred modals.Accounts) (modals.Accounts, error)
	CheckEmailIfExist(username string) bool
	CheckUsernameIfExist(email string) bool
	GetAccountByEmail(email string) (*modals.Accounts, error)
	ForgotPasswordRequestRepo(pwr modals.OTP)
	GetLatestOTP(id int) (modals.OTP, error)
	UpdateOTPStatus(otpID int, isUsed bool)error
}

type InjectLoginDB struct {
	db *gorm.DB
}

func LoginRepoInit(db *gorm.DB) LoginRepository {
	return &InjectLoginDB{db}
}

// Login
func (r *InjectLoginDB) LoginRepo(username string) (modals.Accounts, error) {
	var account modals.Accounts
	var notFoundError = errors.New("user not found")
	err := r.db.Find(&account, "username = ?", username).Error
	if err != nil {
		return modals.Accounts{}, notFoundError
	}
	return account, nil
}

// Register
func (r *InjectLoginDB) RegisterRepo(cred modals.Accounts) (modals.Accounts, error) {
	if err := r.db.Create(&cred).Error; err != nil {
		return modals.Accounts{}, err
	}
	return cred, nil
}

// Check Username if Already Exist
func (r *InjectLoginDB) CheckUsernameIfExist(username string) bool {
	var count int64
	var account modals.Accounts
	r.db.Model(&account).Where("username = ?", username).Count(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}

// Check Email if Already Exist
func (r *InjectLoginDB) CheckEmailIfExist(email string) bool {
	var count int64
	var account modals.Accounts
	r.db.Model(&account).Where("email = ?", email).Count(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}

// Get Account ID using email
func (r *InjectLoginDB) GetAccountByEmail(email string) (*modals.Accounts, error) {
	var account modals.Accounts
	// Use SELECT to only fetch specific columns (optional)
	err := r.db.
		Select("id").
		Where("email = ?", email).
		First(&account).Error

	if err != nil {
		// Record not found
		return nil, err
	}
	return &account, nil
}

// ForgotPassword Request - Save OTP details to DB
func (r *InjectLoginDB) ForgotPasswordRequestRepo(pwr modals.OTP) {
	r.db.Create(&pwr)
}

// Rate Limit
func (r *InjectLoginDB) GetLatestOTP(id int) (modals.OTP, error) {
	var latestRequest modals.OTP
	err := r.db.Where("accounts_id = ?", id).Order("created_at DESC").First(&latestRequest).Error
	if err != nil {
		return modals.OTP{}, err
	}
	return latestRequest, nil
}

//Update OTP Status if already used
func (r *InjectLoginDB) UpdateOTPStatus(otpID int, isUsed bool)error {
	var otp modals.OTP
	err := r.db.Model(otp).Where("id = ?", otpID).Update("used",isUsed).Error
	if err != nil {
		return err
	}
	return nil
}