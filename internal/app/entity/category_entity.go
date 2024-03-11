package entity

type Category struct {
	ID     string  `gorm:"primaryKey"`
	Name   string  `gorm:"not null"`
	Events []Event `gorm:"foreignKey:CategoryID"`
}
