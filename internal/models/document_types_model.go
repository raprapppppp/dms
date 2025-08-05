package models

import (
	"time"
)

type DocumentType struct {
	DocTypeID   int       `json:"docTypeID"   gorm:"column:docTypeID;primaryKey;autoIncrement"`
	DocTypeName string    `json:"docTypeName" gorm:"column:docTypeName;type:varchar(255);not null"`
	DateCreated time.Time `json:"dateCreated" gorm:"column:dateCreated;autoCreateTime"`
	UserID      int       `json:"userID"      gorm:"column:userID"`
	DocID       int       `json:"docID"       gorm:"column:docID"`
	Tag         string    `json:"tag"         gorm:"column:tag;type:varchar(255)"`
}