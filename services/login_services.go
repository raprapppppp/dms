package services

import (
	"dms-api/modals"
	"dms-api/repository"
)
type LoginServices interface {
	LoginService(acc modals.Login) (modals.Login, string)
}

type LoginRepository struct {
	repo repository.LoginRepository
}



