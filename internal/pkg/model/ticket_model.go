package model

import "time"

type TicketUserResponse struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	EventID   string        `json:"event_id"`
	InvoiceID string        `json:"invoice_id"`
	CreatedAt time.Time     `json:"created_at"`
	Event     EventResponse `json:"event"`
}

type TicketOrganizerResponse struct {
	ID        string              `json:"id"`
	UserID    string              `json:"user_id"`
	EventID   string              `json:"event_id"`
	InvoiceID string              `json:"invoice_id"`
	CreatedAt time.Time           `json:"created_at"`
	User      ProfileUserResponse `json:"user"`
}
