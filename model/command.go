package model

import "time"

type CommandStatus string

const (
	Livree    CommandStatus = "Livrée"
	EnAttente CommandStatus = "En attente"
	EnCours   CommandStatus = "En cours"
)

type Command struct {
	ID             uint          `gorm:"primaryKey" json:"id"`
	Status         CommandStatus `gorm:"not null" json:"status"`
	DateCreation   string        `gorm:"not null" json:"date_creation"`
	CommandLine    []CommandLine `json:"command_line"`
	ClientID       uint          `json:"client_id"`
	Client         Client        `gorm:"foreignKey:ClientID" json:"client"`
	LivreurID      uint          `json:"livreur_id"`
	Livreur        User          `gorm:"foreignKey:LivreurID" json:"livreur"`
	DateCompletion string        `json:"date_completion"`
}

func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}