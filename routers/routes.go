package routers

import (
	"simpler-products/controllers"
	"simpler-products/middlewares"
	"simpler-products/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRouter(log *logrus.Logger, ps services.ProductsServiceInterface) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.JSONLoggerMiddleware())
	router.Use(middlewares.GinErrorHandlerMiddleware(log))

	api := router.Group("/api")

	// ping routes
	{
		ping := api.Group("/ping")
		ping.GET("", controllers.Ping())
	}

	// products routes
	{
		products := api.Group("/products")
		products.GET("", controllers.GetAllProducts(ps))
		products.GET("/:id", controllers.GetProductById(ps))
		products.POST("", controllers.AddProduct(ps))
		products.PUT("/:id", controllers.UpdateProduct(ps))
		products.PATCH("/:id", controllers.PatchProduct(ps))
		products.DELETE("/:id", controllers.DeleteProduct(ps))
	}

	return router
}
