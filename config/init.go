package config

import (
	"log"
	"os"

	"simpler-products/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	LogLevel        string
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

	// Set Gin mode based on log level
	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if logLevel == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create services
	productsService := &services.ProductsService{}

	return &Config{
		Port:            port,
		LogLevel:        logLevel,
		ProductsService: productsService,
	}, nil
}
