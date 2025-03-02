package database

import (
	"fmt"
	"kishyassin/Livra-Maroc/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() (*gorm.DB, error) {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Could not load .env file. Falling back to system environment variables.")
	}

	dsn := ""
	if os.Getenv("APP_ENV") == "development" {
		log.Println("Running in development environment.")
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_DATABASE"),
		)
	} else {
		log.Println("Running in production environment.")
		dsn = "root:tETjVsheYljqOLKpimBxZJIVKfqSDuSr@tcp(nozomi.proxy.rlwy.net:34335)/railway?charset=utf8mb4&parseTime=True&loc=Local"
	}
	dialector := mysql.Open(dsn)

	// Connect to the database
	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database. DSN: %s, Error: %v", dsn, err)
		return nil, err
	}

	log.Println("Database connection established successfully.")

	// Run migrations
	if err := runMigrations(DB); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
		return nil, err
	}

	log.Println("Database migration completed successfully.")
	return DB, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Client{},
		&model.Command{},
		&model.CommandLine{},
		&model.User{},
		&model.Product{},
	)
}