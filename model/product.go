package model

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `gorm:"not null" json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}
