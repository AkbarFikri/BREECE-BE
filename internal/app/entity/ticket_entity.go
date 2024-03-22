package entity

import "time"

type Ticket struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	InvoiceID string    `gorm:"not null"`
	EventID   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	Event     Event
	User      User
}
