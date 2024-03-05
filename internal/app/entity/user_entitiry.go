package entity

type User struct {
	ID       string `gorm:"primaryKey"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Username string `gorm:"not null"`
}
