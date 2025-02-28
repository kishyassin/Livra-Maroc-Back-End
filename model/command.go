package model

type Command struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Status       string        `gorm:"not null" json:"status"`
	DateCreation string        `gorm:"not null" json:"date_creation"`
	CommandLine  []CommandLine `json:"command_line"`
	ClientID     uint          `json:"client_id"`
	Client       Client        `gorm:"foreignKey:ClientID" json:"client"`
}
