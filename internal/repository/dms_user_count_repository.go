package repository

import (
	"dms-api/internal/models"

	"gorm.io/gorm"
)

// All Available Repository in DMSUsersRepository
type DMSUsersRepository interface {
	GetAllDMSActiveUsersCount() (int, error)
	GetAllDMSNonActiveUsersCount() (int, error)
	GetAllDMSActiveUsersBranchCount(branchID uint) (int, error)
	GetAllDMSNonActiveUsersBranchCount(branchID uint) (int, error)
}

type DMSUsersDBGorm struct {
	db *gorm.DB
}

// Repository Initializer
func DMSUserRepositorysInit(db *gorm.DB) DMSUsersRepository {
	return &DMSUsersDBGorm{db}
}

func (d *DMSUsersDBGorm) GetAllDMSActiveUsersCount() (int, error) {
	var count int64
	err := d.db.Model(&models.User{}).Where("status = ?", "active").Count(&count).Error
	return int(count), err
}

func (d *DMSUsersDBGorm) GetAllDMSNonActiveUsersCount() (int, error) {
	var count int64
	err := d.db.Model(&models.User{}).Where("status = ?", "non-active").Count(&count).Error
	return int(count), err
}

func (d *DMSUsersDBGorm) GetAllDMSActiveUsersBranchCount(branchID uint) (int, error) {
	var count int64
	err := d.db.Model(&models.User{}).
		Where("status = ? AND branch_id = ?", "active", branchID).
		Count(&count).Error
	return int(count), err
}

func (d *DMSUsersDBGorm) GetAllDMSNonActiveUsersBranchCount(branchID uint) (int, error) {
	var count int64
	err := d.db.Model(&models.User{}).
		Where("status = ? AND branch_id = ?", "non-active", branchID).
		Count(&count).Error
	return int(count), err
}
