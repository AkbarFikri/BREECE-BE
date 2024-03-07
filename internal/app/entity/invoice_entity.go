package entity

import "time"

type Invoice struct {
	ID        string
	UserID    string
	Status    string
	Snap      string
	Ticket    Ticket
	CreatedAt time.Time
}
