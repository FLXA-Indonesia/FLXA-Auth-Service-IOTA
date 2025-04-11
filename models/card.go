package models

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	CardPhoneNumber string `gorm:"unique"`
	CardStatus      string
	OperatorID      uuid.UUID
	UserID          uint
	CardDateAdded   time.Time
	CardID          uint `gorm:"primaryKey"`
	User            User `gorm:"foreignKey:UserID"`
}

func (Card) TableName() string {
	return "Card"
}
