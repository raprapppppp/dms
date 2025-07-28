package repository

import (
	"dms-api/modals"
	"errors"

	"gorm.io/gorm"
)

type LoginRepository interface {
	Login(username string) (modals.Login, error)
}

type LoginInjectDB struct {
	db *gorm.DB
}

func LoginRepoInit(db *gorm.DB) LoginRepository{
	return &LoginInjectDB{db}
}

func(r *LoginInjectDB) Login(username string)(modals.Login, error){
var account modals.Login
var notFoundError = errors.New("user not found")
err := r.db.Find(&account, "username = ", username).Error
if err != nil {
	return modals.Login{}, notFoundError
}
return account,nil
}

