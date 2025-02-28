package model
type User struct {
	ID       uint   `gorm:"primaryKey"`
	FirstName string `gorm:"not null"`
	LastName string `gorm:"not null"`
	Password string
	Email    string `gorm:"unique"`
	Role    string `gorm:"default:'student'"`
}
