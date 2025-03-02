package model

type CommandLine struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	ProductID uint `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int  `json:"quantity"`
	CommandID uint `json:"command_id"`
}
