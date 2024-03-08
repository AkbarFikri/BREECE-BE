package entity

import (
	"time"

	"github.com/lib/pq"

)

type Event struct {
	ID          string `gorm:"primaryKey"`
	CategoryID  string
	Title       string
	Description string
	Tempat      string
	Speakers    pq.StringArray
	Date        time.Time
	StartAt     time.Time
	Link        string
	Price       uint32
	TicketQty   uint16
	OrganizeBy  string
	Tickets     []Ticket `gorm:"foreignKey:EventID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
