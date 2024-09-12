package main

import (
	"fmt"
	"log"
	"simpler-products/config"
	"simpler-products/routers"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	// close Database connection when app terminates
	defer cfg.DB.Close()
	defer fmt.Println("Closing Database connection")

	router := routers.NewRouter(cfg.ProductsService)

	log.Fatal(router.Run(":" + cfg.Port))
}
