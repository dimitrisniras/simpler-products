package routers

import (
	"simpler-products/config"
	"simpler-products/controllers"
	"simpler-products/middlewares"
	"simpler-products/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRouter(servs config.ServiceContainer, log *logrus.Logger) *gin.Engine {
	router := gin.New()

	// add middlewares
	router.Use(gin.Recovery())
	router.Use(middlewares.SecurityHeaders())
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.ResponseFormatter(log))
	router.Use(middlewares.JSONLoggerMiddleware())

	// generic /api endpoint
	api := router.Group("/api")

	// /ping routes
	{
		ping := api.Group("/ping")
		ping.GET("", controllers.Ping())
	}

	// /products routes
	productsService, ok := servs.(services.ProductsServiceInterface)
	if !ok {
		log.Fatal("ProductsServiceInterface not found in services")
	}

	{
		products := api.Group("/products")

		// use auth middleware
		products.Use(middlewares.JWTAuthMiddleware())

		products.GET("", controllers.GetAllProducts(productsService))
		products.GET("/:id", controllers.GetProductById(productsService))
		products.POST("", controllers.AddProduct(productsService))
		products.PUT("/:id", controllers.UpdateProduct(productsService))
		products.DELETE("/:id", controllers.DeleteProduct(productsService))
	}

	return router
}
