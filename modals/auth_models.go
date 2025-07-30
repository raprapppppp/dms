package modals

import "time"

type Login struct {
	Username string
	Password string
}

type Accounts struct {
	ID       uint   `gorm:"primaryKey;autoIncreament"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type Forgot struct {
	Email string `json:"email"`
}

type OTP struct {
	ID         uint      `gorm:"primaryKey"`
	AccountsID uint      `gorm:"not null;index"`
	OTPCode    string    `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	Used       bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"default:current_timestamp"`
	Accounts   Accounts  `gorm:"foreignKey:AccountsID;constraint:OnDelete:CASCADE"`
}

type VerifyOTP struct {
	Identifier string
	Otp string
}
