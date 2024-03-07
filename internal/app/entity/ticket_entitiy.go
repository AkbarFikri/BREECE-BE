package entity

import "time"

type Ticket struct {
	ID        string
	UserID    string
	InvoiceID string
	EventID   string
	CreatedAt time.Time
}
