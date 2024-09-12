package routers

import (
	"encoding/json"
	"net/http"
	"simpler-products/controllers"
	"simpler-products/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func JSONLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()
			log["user_agent"] = params.Request.UserAgent()
			log["error"] = params.ErrorMessage

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}

func ErrorHandlerMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// Log the error
			log.Println(c.Errors.ByType(gin.ErrorTypePrivate).String())

			// Send a generic error response to the client
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
	}
}

func NewRouter(log *logrus.Logger, ps services.ProductsServiceInterface) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(JSONLoggerMiddleware())
	router.Use(ErrorHandlerMiddleware(log))

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
