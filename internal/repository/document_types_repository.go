package repository

import (
	"dms-api/internal/models"

	"gorm.io/gorm"
)

//All Available Repository in DT
type DocumentTypesRepository interface {
	GetAllDocumentTypes() ([]models.DocumentTypes, error)
}

type DBGorm struct {
	db *gorm.DB
}
//Repository Initializer
func DocumentTypesRepositoryInit(db *gorm.DB) DocumentTypesRepository {
	return &DBGorm{db}
}

func(d *DBGorm) GetAllDocumentTypes() ([]models.DocumentTypes, error){
	var documentTypes []models.DocumentTypes

	if err := d.db.Model(&models.DocumentTypes{}).Find(&documentTypes).Error; err != nil {
		return nil, err
	}
	return documentTypes, nil
}