package routers

import (
	"simpler-products/controllers"
	"simpler-products/services"

	"github.com/gin-gonic/gin"
)

func NewRouter(ps services.ProductsServiceInterface) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

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
