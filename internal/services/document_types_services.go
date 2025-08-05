package services

import (
	"dms-api/internal/models"
	"dms-api/internal/repository"
)

//All Available Services in DT
type DocumentTypesServices interface {
	GetAllDocumentTypes() ([]models.DocumentTypes, error)
}

type InjectDocumentTypesRepository struct {
	repository repository.DocumentTypesRepository
}
//Service Initializer
func DocumentTypesServicesInit(repo repository.DocumentTypesRepository) DocumentTypesServices {
	return &InjectDocumentTypesRepository{repo}
}

//Get All DT
func(s *InjectDocumentTypesRepository) GetAllDocumentTypes()([]models.DocumentTypes, error){
	return s.repository.GetAllDocumentTypes()
}