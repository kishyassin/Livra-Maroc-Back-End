package main

import (
	"kishyassin/Livra-Maroc/database"
	"kishyassin/Livra-Maroc/router"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize the database connection
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Set up the routes with the database instance
	router.SetupRoutes(app, db)

	// Start the server
	port := os.Getenv("APP_PORT") // Railway assigns this dynamically
	if port == "" {
		port = "3000" // Default fallback
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
