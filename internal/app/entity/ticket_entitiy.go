package entity

import "time"

type Ticket struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	InvoiceID string
	EventID   string
	CreatedAt time.Time
}
