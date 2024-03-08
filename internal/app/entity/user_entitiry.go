package entity

import "time"

type User struct {
	ID                string `gorm:"primaryKey"`
	Email             string `gorm:"not null;unique"`
	Password          string `gorm:"not null"`
	Nim               uint64
	FullName          string
	Prodi             string
	Universitas       string
	IsEmailVerified   bool
	EmailVerifiedAt   time.Time
	IsProfileVerified bool
	IsAdmin           bool
	IsOrganizer       bool
	ID_Url            string
	Invoices          []Invoice `gorm:"foreignKey:UserID"`
	Events            []Event   `gorm:"foreignKey:OrganizeBy"`
	Ticekts           []Ticket  `gorm:"foreignKey:UserID"`
}
