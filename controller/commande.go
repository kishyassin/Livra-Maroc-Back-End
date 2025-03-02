package controller

import (
	"kishyassin/Livra-Maroc/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetNonCompletedCommandes fetches all commandes of a livreur  that are not completed status = "en cours"
func GetNonCompletedCommandes(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		livreurId, ok := c.Locals("userID").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}
		var commandes []model.Command
		db.Preload("Client").Preload("CommandLine").Where("livreur_id = ? AND status != ?", livreurId, "Livrée").Find(&commandes)

		type Response struct {
			ID            uint                `json:"id"`
			Destination   string              `json:"destination"`
			Status        string              `json:"status"`
			NumberOfLines int                 `json:"number_of_lines"`
			ClientName    string              `json:"client_name"`
			CommandLines  []model.CommandLine `json:"command_lines"`
		}

		var response []Response
		for _, commande := range commandes {
			response = append(response, Response{
				ID:            commande.ID,
				Destination:   commande.Client.Location,
				Status:        string(commande.Status),
				NumberOfLines: len(commande.CommandLine),
				ClientName:    commande.Client.FirstName + " " + commande.Client.LastName,
				CommandLines:  commande.CommandLine,
			})
		}

		return c.JSON(response)
	}
}

// GetCompletedTodayCommandes fetches all commandes of a livreur that are completed status = "Livrée" and completed today
func GetCompletedTodayCommandes(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		livreurId, ok := c.Locals("userID").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		var commandes []model.Command
		today := time.Now().Format("2006-01-02") // Format YYYY-MM-DD

		// Requête corrigée avec la date du jour
		db.Preload("Client").Preload("CommandLine").
			Where("livreur_id = ? AND status = ? AND DATE(date_completion) = ?",
				livreurId, "Livrée", today).
			Find(&commandes)

		type Response struct {
			ID            uint                `json:"id"`
			Destination   string              `json:"destination"`
			Status        string              `json:"status"`
			NumberOfLines int                 `json:"number_of_lines"`
			ClientName    string              `json:"client_name"`
			CommandLines  []model.CommandLine `json:"command_lines"`
		}

		var response []Response
		for _, commande := range commandes {
			response = append(response, Response{
				ID:            commande.ID,
				Destination:   commande.Client.Location,
				Status:        string(commande.Status),
				NumberOfLines: len(commande.CommandLine),
				ClientName:    commande.Client.FirstName + " " + commande.Client.LastName,
				CommandLines:  commande.CommandLine,
			})
		}

		return c.JSON(response)
	}
}

// GetCommandesSummary fetches a summary of commandes for a livreur
func GetCommandesSummary(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		livreurId, ok := c.Locals("userID").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		var completedCommandes []model.Command
		var nonCompletedCommandes []model.Command
		today := time.Now().Format("2006-01-02") // Format YYYY-MM-DD

		// Fetch completed commandes
		db.Preload("Client").Preload("CommandLine.Product").
			Where("livreur_id = ? AND status = ? AND DATE(date_completion) = ?",
				livreurId, "Livrée", today).
			Find(&completedCommandes)

		// Fetch non-completed commandes
		db.Preload("Client").Preload("CommandLine.Product").
			Where("livreur_id = ? AND status != ?", livreurId, "Livrée").
			Find(&nonCompletedCommandes)

		type CommandLineResponse struct {
			ProductName string `json:"product_name"`
			Quantity    int    `json:"quantity"`
		}

		type Response struct {
			ID            uint                  `json:"id"`
			Destination   string                `json:"destination"`
			Status        string                `json:"status"`
			NumberOfLines int                   `json:"number_of_lines"`
			ClientName    string                `json:"client_name"`
			CommandLines  []CommandLineResponse `json:"command_lines"`
		}

		var completedResponse = []Response{}
		for _, commande := range completedCommandes {
			var commandLines []CommandLineResponse
			for _, line := range commande.CommandLine {
				commandLines = append(commandLines, CommandLineResponse{
					ProductName: line.Product.Name,
					Quantity:    line.Quantity,
				})
			}
			completedResponse = append(completedResponse, Response{
				ID:            commande.ID,
				Destination:   commande.Client.Location,
				Status:        string(commande.Status),
				NumberOfLines: len(commande.CommandLine),
				ClientName:    commande.Client.FirstName + " " + commande.Client.LastName,
				CommandLines:  commandLines,
			})
		}

		var nonCompletedResponse = []Response{}
		for _, commande := range nonCompletedCommandes {
			var commandLines []CommandLineResponse
			for _, line := range commande.CommandLine {
				commandLines = append(commandLines, CommandLineResponse{
					ProductName: line.Product.Name,
					Quantity:    line.Quantity,
				})
			}
			nonCompletedResponse = append(nonCompletedResponse, Response{
				ID:            commande.ID,
				Destination:   commande.Client.Location,
				Status:        string(commande.Status),
				NumberOfLines: len(commande.CommandLine),
				ClientName:    commande.Client.FirstName + " " + commande.Client.LastName,
				CommandLines:  commandLines,
			})
		}

		return c.JSON(fiber.Map{
			"totalDeliveries":          len(completedCommandes) + len(nonCompletedCommandes),
			"completedDeliveriesCount": len(completedCommandes),
			"remainingDeliveriesCount": len(nonCompletedCommandes),
			"completedDeliveries":      completedResponse,
			"remainingDeliveries":      nonCompletedResponse,
		})
	}
}


// UpdateCommandeStatus updates the status of a commande
func UpdateCommandeStatus(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request struct {
			CommandID uint   `json:"command_id"`
			Status    string `json:"status"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Check if the status is valid
		if request.Status != "Livrée" && request.Status != "En cours" && request.Status != "En attente" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status",
			})
		}

		// Check if the commande exists
		var commande model.Command
		if err := db.Where("id = ?", request.CommandID).First(&commande).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Commande not found",
			})
		}

		// Update the status
		commande.Status = model.CommandStatus(request.Status)
		// Update the date of completion if the status is "Livrée"
		if request.Status == "Livrée" {
			commande.DateCompletion = model.GetTodayDate()
		}else{
			commande.DateCompletion = ""
		}
		if err := db.Save(&commande).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not update status",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Status updated successfully",
		})
	}
}