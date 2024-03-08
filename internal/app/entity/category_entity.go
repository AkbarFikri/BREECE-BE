package entity

type Category struct {
	ID     string `gorm:"primaryKey"`
	Name   string
	Events []Event `gorm:"foreignKey:CategoryID"`
}
