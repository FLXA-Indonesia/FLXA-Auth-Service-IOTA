package models

type User struct {
	UserID       uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	Email        string ``
	SecretString string `gorm:"not null"`
	ProfilePhoto string
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return "User"
}
