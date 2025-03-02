package router

import (
	"kishyassin/Livra-Maroc/controller"
	"kishyassin/Livra-Maroc/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	//health
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Supreme the KING \n Don't mess with me")
	})

	// Auth Routes
	auth := app.Group("/auth")
	auth.Post("/login", controller.Login(db))
	auth.Post("/refresh", controller.RefreshToken())

	livreur := app.Group("/livreur")
	livreur.Get("non-completed-commandes", middleware.Authenticated(), controller.GetNonCompletedCommandes(db))
	livreur.Get("completed-today-commandes", middleware.Authenticated(), controller.GetCompletedTodayCommandes(db))
	livreur.Get("commandes-summary", middleware.Authenticated(), controller.GetCommandesSummary(db))
	livreur.Patch("update-commande-status", middleware.Authenticated(), controller.UpdateCommandeStatus(db))
}
