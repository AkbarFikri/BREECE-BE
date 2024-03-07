package entity

import (
	"time"

	"github.com/lib/pq"
)

type Event struct {
	ID           string
	CategoryID   string
	CategoryName string
	Title        string
	Description  string
	Tempat       string
	Speakers     pq.StringArray
	Date         time.Time
	StartAt      time.Time
	Link         string
	Price        uint32
	TicketQty    uint16
	Tickets      []Ticket
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
