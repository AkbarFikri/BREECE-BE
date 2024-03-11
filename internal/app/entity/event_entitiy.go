package entity

import (
	"time"

	"github.com/lib/pq"

)

type Event struct {
	ID           string         `gorm:"primaryKey"`
	CategoryID   string         `gorm:"not null"`
	Title        string         `gorm:"not null"`
	Description  string         `gorm:"not null"`
	Tempat       string         `gorm:"not null"`
	Speakers     pq.StringArray `gorm:"not null;type:text[]"`
	SpeakersRole pq.StringArray `gorm:"not null;type:text[]"`
	Date         time.Time      `gorm:"not null"`
	StartAt      time.Time      `gorm:"not null"`
	Link         string         `gorm:"not null"`
	Price        uint32         `gorm:"not null"`
	TicketQty    uint16         `gorm:"not null"`
	OrganizeBy   string         `gorm:"not null"`
	Tickets      []Ticket       `gorm:"foreignKey:EventID"`
	Invoices     []Invoice      `gorm:"foreignKey:"EventID"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
}
