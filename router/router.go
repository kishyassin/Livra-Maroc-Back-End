package router

import (
	"kishyassin/Livra-Maroc/controller"

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
}
