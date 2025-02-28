package model

type CommandLine struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
	CommandID uint `json:"command_id"`
}
