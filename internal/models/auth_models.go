package models

type Accounts struct {
	ID       uint   `gorm:"primaryKey;autoIncreament"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}
