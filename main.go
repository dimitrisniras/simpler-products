package main

import (
	"log"
	"simpler-products/config"
	"simpler-products/routers"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	router := routers.NewRouter(cfg.ProductsService)

	log.Fatal(router.Run(":" + cfg.Port))
}
