package repository

import (
	"dms-api/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AuthRepository interface {
	LoginRepo(username string) (models.User, error)
	//RegisterRepo(cred models.Accounts) (models.Accounts, error)
	CheckStaffIdIfExist(staffId string) bool
	CheckUsernameIfExist(username string) bool
	GetUserIDByStaffID(staffId string) (*models.User, error)
	ForgotPasswordRequestRepo(pwr models.OTP)
	GetLatestOTP(id int) (models.OTP, error)
	UpdateOTPStatus(otpID int, isUsed bool) error
	UpdatePasswordRepo(id int, newPassword string) error
}

type InjectLoginDB struct {
	db *gorm.DB
}

func AuthRepoInit(db *gorm.DB) AuthRepository {
	return &InjectLoginDB{db}
}

// Login 
func (r *InjectLoginDB) LoginRepo(username string) (models.User, error) {
	var user models.User
	var notFoundError = errors.New("user not found")
	err := r.db.Find(&user, "user_name = ?", username).Error
	if err != nil {
		return models.User{}, notFoundError
	}
	return user, nil
}
// Register
/* func (r *InjectLoginDB) RegisterRepo(cred models.Accounts) (models.Accounts, error) {
	if err := r.db.Create(&cred).Error; err != nil {
		return models.Accounts{}, err
	}
	return cred, nil
} */

// Check Username if Already Exist
func (r *InjectLoginDB) CheckUsernameIfExist(username string) bool {
	var count int64
	var user models.User
	r.db.Model(&user).Where("user_name = ?", username).Count(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}

// Check Staff ID if Already Exist 
func (r *InjectLoginDB) CheckStaffIdIfExist(staffId string) bool {
	var count int64
	var user models.User
	r.db.Model(&user).Where("staff_id = ?", staffId).Count(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}

// Get User ID using Staff ID
func (r *InjectLoginDB) GetUserIDByStaffID(staffId string) (*models.User, error) {
	var user models.User
	// Use SELECT to only fetch specific columns (optional)
	result := r.db.Select("user_id").Where("staff_id = ?", staffId).First(&user)
	if result.Error != nil {
		return &models.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return &models.User{}, fmt.Errorf("no record found with id %s", staffId)
	}
	return &user, nil
}

// ForgotPassword Request - Save OTP details to DB
func (r *InjectLoginDB) ForgotPasswordRequestRepo(pwr models.OTP) {
	r.db.Create(&pwr)
}

// Rate Limit
func (r *InjectLoginDB) GetLatestOTP(id int) (models.OTP, error) {
	var latestRequest models.OTP
	result := r.db.Where("user_id = ?", id).Order("created_at DESC").First(&latestRequest)
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
func (r *InjectLoginDB) UpdatePasswordRepo(userId int, newPassword string)error{
	var user models.User
	result := r.db.Model(&user).Where("user_id = ?", userId).Update("user_pass",newPassword)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no record found with id %v", userId)
	}
	return nil
}
