package main

import (
	"os"
	"simpler-products/config"
	"simpler-products/routers"

	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	cfg, err := config.Init(log)
	if err != nil {
		log.Fatal(err)
	}

	// close Database connection when app terminates
	defer cfg.DB.Close()
	defer log.Debug("Closing Database connection")

	router := routers.NewRouter(cfg.Services, log)

	log.Printf("Server listening on :%s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
