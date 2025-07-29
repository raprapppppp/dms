package modals

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
	Email string
}