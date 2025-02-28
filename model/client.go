package model

type Client struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Phone     string `gorm:"not null" json:"phone"`
	Location  string `gorm:"not null" json:"location"`
}
