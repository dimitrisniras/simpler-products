package config

import (
	"database/sql"
	"log"
	"os"

	"simpler-products/database"
	"simpler-products/services"

	_ "github.com/go-sql-driver/mysql" // MySQL driver

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	LogLevel        string
	DB              *sql.DB
	ProductsService services.ProductsServiceInterface
}

func Init() (*Config, error) {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	port := os.Getenv("PORT")
	logLevel := os.Getenv("LOG_LEVEL")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Set Gin mode based on log level
	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if logLevel == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Database setup
	db, err := database.Init(dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		return nil, err
	}

	// Create services
	productsService := &services.ProductsService{}

	return &Config{
		Port:            port,
		LogLevel:        logLevel,
		DB:              db,
		ProductsService: productsService,
	}, nil
}
