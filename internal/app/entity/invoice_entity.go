package entity

import "time"

type Invoice struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	EventID   string    `gorm:"not null"`
	Amount    int64     `gorm:not null`
	Status    string    `gorm:"not null"`
	Snap      string    `gorm:"not null"`
	Ticket    Ticket    `gorm:"foreignKey:InvoiceID"`
	CreatedAt time.Time `gorm:"not null"`
}
