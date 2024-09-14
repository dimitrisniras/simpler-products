package routers

import (
	"os"
	"simpler-products/config"
	"simpler-products/controllers"
	v1Controllers "simpler-products/controllers/v1"
	"simpler-products/middlewares"
	"simpler-products/services"

	"github.com/dvwright/xss-mw"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRouter(servs config.ServiceContainer, log *logrus.Logger) *gin.Engine {
	router := gin.New()

	// add middlewares
	router.Use(gin.Recovery())

	// sanitize input for XSS protection
	var xssMdlwr xss.XssMw
	router.Use(xssMdlwr.RemoveXss())

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

	{
		// v1 routes
		v1Routes := api.Group("/v1")

		// /products routes
		{
			productsService, ok := servs.(services.ProductsServiceInterface)
			if !ok {
				log.Fatal("ProductsServiceInterface not found in services")
			}

			products := v1Routes.Group("/products")

			authEnabled := os.Getenv("AUTH_ENABLED")
			if authEnabled == "true" {
				// use auth middleware
				products.Use(middlewares.JWTAuthMiddleware())
			}

			products.GET("", v1Controllers.GetAllProducts(productsService))
			products.GET("/:id", v1Controllers.GetProductById(productsService))
			products.POST("", v1Controllers.AddProduct(productsService))
			products.PUT("/:id", v1Controllers.UpdateProduct(productsService))
			products.DELETE("/:id", v1Controllers.DeleteProduct(productsService))

		}
	}

	return router
}
