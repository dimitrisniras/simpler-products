package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"simpler-products/config"
	"simpler-products/routers"
	"syscall"
	"time"

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

	router := routers.NewRouter(cfg.Services, log)

	// Create a server instance
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	<-quit
	log.Println("Shutting down server...")

	// close Database connection when app terminates
	defer cfg.DB.Close()
	defer log.Debug("Closing Database connection")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
