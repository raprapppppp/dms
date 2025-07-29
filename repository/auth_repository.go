package repository

import (
	"dms-api/modals"
	"errors"

	"gorm.io/gorm"
)

type LoginRepository interface {
	LoginRepo(username string) (modals.Accounts, error)
	RegisterRepo(cred modals.Accounts) (modals.Accounts,error)
	CheckEmailIfExist(username string) bool
	CheckUsernameIfExist(email string) bool
}

type InjectLoginDB struct {
	db *gorm.DB
}

func LoginRepoInit(db *gorm.DB) LoginRepository{
	return &InjectLoginDB{db}
}
//Login
func(r *InjectLoginDB) LoginRepo(username string)(modals.Accounts, error){
	var account modals.Accounts
	var notFoundError = errors.New("user not found")
	err := r.db.Find(&account, "username = ?", username).Error
	if err != nil {
		return modals.Accounts{}, notFoundError
	}
	return account,nil
}
//Register
func(r *InjectLoginDB) RegisterRepo(cred modals.Accounts) (modals.Accounts,error){
	if err := r.db.Create(&cred).Error; err != nil {
		return modals.Accounts{}, err
	}
	return cred, nil
}
//Check Username if Already Exist
func(r *InjectLoginDB) CheckUsernameIfExist(username string) bool {
	var count int64
	var account modals.Accounts
	r.db.Model(&account).Where("username = ?", username).Count(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}
//Check Email if Already Exist
func(r *InjectLoginDB) CheckEmailIfExist(email string) bool {
	var count int64
	var account modals.Accounts
	r.db.Model(&account).Where("email = ?", email).Count(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}



