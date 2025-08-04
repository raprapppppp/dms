package models

import "time"

type User struct {
	UserID         uint      `gorm:"primaryKey;autoIncrement" json:"user_id"`
	OldEmpID       string    `json:"old_emp_id"`
	MapCID         string    `json:"map_cid"`
	UserName       string    `gorm:"uniqueIndex" json:"user_name"` // Unique
	UserPass       string    `json:"user_pass"`
	RoleID         uint      `json:"role_id"`
	RoleName       string    `json:"role_name"`
	StatusID       uint      `json:"status_id"`
	StatusName     string    `json:"status_name"`
	IsOnline       bool      `json:"is_online"`
	FailedAttempts int       `json:"failed_attempts"`
	DateLocked     time.Time `json:"date_locked"`
	LastLogin      time.Time `json:"last_login"`
	ExpirationDate time.Time `json:"expiration_date"`
	DateCreated    time.Time `gorm:"autoCreateTime" json:"date_created"`
	DaysRemaining  int       `json:"days_remaining"`
	StaffID        string    `gorm:"uniqueIndex" json:"staff_id"` // Unique
	UserPassword   string    `json:"user_password"`
	LastName       string    `json:"last_name"`
	FirstName      string    `json:"first_name"`
	MiddleName     string    `json:"middle_name"`
	FullName       string    `json:"full_name"`
	BirthDate      string    `json:"birth_date"`
	EmailAddress   string    `json:"email_address"`
	MobileNumber   string    `json:"mobile_number"`
	UserToken      string    `json:"user_token"`
	UsedToken      bool      `json:"used_token"`
	Tag            string    `json:"tag"`
	Filter         string    `json:"filter"`
	UID            uint      `json:"uid"`
	InstiID        uint      `json:"insti_id"`
	OfficeID       string    `json:"office_id"`
}