package entity

import "time"

type User struct {
	ID                string `gorm:"primaryKey"`
	Email             string `gorm:"not null;unique"`
	Password          string
	NimNik            string
	FullName          string
	Prodi             string
	Universitas       string
	IsEmailVerified   bool      `gorm:"not null"`
	EmailVerifiedAt   time.Time `gorm:"not null"`
	IsProfileVerified bool      `gorm:"not null"`
	IsAdmin           bool      `gorm:"not null"`
	IsOrganizer       bool      `gorm:"not null"`
	IsBrawijaya       bool      `gorm:"not null;default:false"`
	IDUrl             string    `gorm:"not null"`
	Invoices          []Invoice `gorm:"foreignKey:UserID"`
	Events            []Event   `gorm:"foreignKey:OrganizeBy"`
	Ticekts           []Ticket  `gorm:"foreignKey:UserID"`
}
