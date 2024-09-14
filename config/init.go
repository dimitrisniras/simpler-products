package config

import (
	"database/sql"
	"os"

	"simpler-products/database"
	"simpler-products/services"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Define a generic interface for any service that might be used by the application
type ServiceContainer interface{}

type Config struct {
	Port     string
	DB       *sql.DB
	Services ServiceContainer
	Log      *logrus.Logger
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
	switch logLevel {
	case "trace":
		gin.SetMode(gin.TestMode)
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		gin.SetMode(gin.DebugMode)
		log.SetLevel(logrus.DebugLevel)
	case "info":
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.WarnLevel)
	case "error":
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.ErrorLevel)
	case "release":
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.InfoLevel)
	default:
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.InfoLevel)
	}

	// Database setup
	db, err := database.Init(log, dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		return nil, err
	}

	// Create services and store them in a struct implementing ServiceContainer
	services := struct {
		services.ProductsServiceInterface
	}{
		&services.ProductsService{
			DB:  db,
			Log: log,
		},
	}

	return &Config{
		Port:     port,
		DB:       db,
		Services: services,
		Log:      log,
	}, nil
}
