package entity

import "time"

type Invoice struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	Status    string
	Snap      string
	Ticket    Ticket `gorm:"foreignKey:InvoiceID"`
	CreatedAt time.Time
}
