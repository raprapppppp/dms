package services

import "dms-api/internal/repository"

//All Available Services in DMSUsersServices
type DMSUsersServices interface {
	GetAllDocumentTypes() (int, error)
}

type InjectDMSUsersRepository struct {
	repository repository.DMSUsersRepository
}

//Service Initializer
func DMSUserServicesInit(repo repository.DMSUsersRepository) DMSUsersServices {
	return &InjectDMSUsersRepository{repo}
}

//
func (s *InjectDMSUsersRepository) GetAllDocumentTypes() (int, error) {
	return s.repository.GetAllDMSActiveUsersCount()
}
