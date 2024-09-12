package config

import (
	"database/sql"
	"os"

	"simpler-products/database"
	"simpler-products/services"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	LogLevel        string
	DB              *sql.DB
	ProductsService services.ProductsServiceInterface
	Log             *logrus.Logger
}

func Init(log *logrus.Logger) (*Config, error) {
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

	// Set Gin mode and stdout logs based on log level
	if logLevel == "trace" {
		gin.SetMode(gin.TestMode)
		log.SetLevel(logrus.TraceLevel)
	} else if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
		log.SetLevel(logrus.DebugLevel)
	} else if logLevel == "info" {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.InfoLevel)
	} else if logLevel == "warn" {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.WarnLevel)
	} else if logLevel == "error" {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.ErrorLevel)
	} else if logLevel == "release" {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.InfoLevel)
	}

	// Database setup
	db, err := database.Init(log, dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		return nil, err
	}

	// Create services
	productsService := &services.ProductsService{
		DB:  db,
		Log: log,
	}

	return &Config{
		Port:            port,
		LogLevel:        logLevel,
		DB:              db,
		ProductsService: productsService,
		Log:             log,
	}, nil
}
