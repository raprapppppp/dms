package repository

import (
	"dms-api/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type LoginRepository interface {
	LoginRepo(username string) (models.Accounts, error)
	RegisterRepo(cred models.Accounts) (models.Accounts, error)
	CheckEmailIfExist(username string) bool
	CheckUsernameIfExist(email string) bool
	GetAccountByEmail(email string) (*models.Accounts, error)
	ForgotPasswordRequestRepo(pwr models.OTP)
	GetLatestOTP(id int) (models.OTP, error)
	UpdateOTPStatus(otpID int, isUsed bool) error
	UpdatePasswordRepo(id int, newPassword string) error
}

type InjectLoginDB struct {
	db *gorm.DB
}

func LoginRepoInit(db *gorm.DB) LoginRepository {
	return &InjectLoginDB{db}
}

// Login
func (r *InjectLoginDB) LoginRepo(username string) (models.Accounts, error) {
	var account models.Accounts
	var notFoundError = errors.New("user not found")
	err := r.db.Find(&account, "username = ?", username).Error
	if err != nil {
		return models.Accounts{}, notFoundError
	}
	return account, nil
}

// Register
func (r *InjectLoginDB) RegisterRepo(cred models.Accounts) (models.Accounts, error) {
	if err := r.db.Create(&cred).Error; err != nil {
		return models.Accounts{}, err
	}
	return cred, nil
}

// Check Username if Already Exist
func (r *InjectLoginDB) CheckUsernameIfExist(username string) bool {
	var count int64
	var account models.Accounts
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
	var account models.Accounts
	r.db.Model(&account).Where("email = ?", email).Count(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}

// Get Account ID using email
func (r *InjectLoginDB) GetAccountByEmail(email string) (*models.Accounts, error) {
	var account models.Accounts
	// Use SELECT to only fetch specific columns (optional)
	result := r.db.Select("id").Where("email = ?", email).First(&account)
	if result.Error != nil {
		return &models.Accounts{}, result.Error
	}
	if result.RowsAffected == 0 {
		return &models.Accounts{}, fmt.Errorf("no record found with id %s", email)
	}
	return &account, nil
}

// ForgotPassword Request - Save OTP details to DB
func (r *InjectLoginDB) ForgotPasswordRequestRepo(pwr models.OTP) {
	r.db.Create(&pwr)
}

// Rate Limit
func (r *InjectLoginDB) GetLatestOTP(id int) (models.OTP, error) {
	var latestRequest models.OTP
	result := r.db.Where("accounts_id = ?", id).Order("created_at DESC").First(&latestRequest)
	if result.Error != nil {
		return models.OTP{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.OTP{}, fmt.Errorf("no record found with id %v", id)
	}
	return latestRequest, nil
}

// Update OTP Status if already used
func (r *InjectLoginDB) UpdateOTPStatus(otpID int, isUsed bool) error {
	var otp models.OTP
	err := r.db.Model(&otp).Where("id = ?", otpID).Update("used", isUsed).Error
	if err != nil {
		return err
	}
	return nil
}

//Update Password in DB
func (r *InjectLoginDB) UpdatePasswordRepo(id int, newPassword string)error{
	var account models.Accounts
	result := r.db.Model(&account).Where("id = ?", id).Update("password",newPassword)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no record found with id %v", id)
	}
	return nil
}
