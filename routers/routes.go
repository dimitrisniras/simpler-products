package routers

import (
	"encoding/json"
	"net/http"
	"simpler-products/controllers"
	"simpler-products/services"
	"simpler-products/validators"

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
			log.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())

			// Construct a structured error response with an array of errors
			var errorResponse struct {
				Status int                 `json:"status"`
				Errors []map[string]string `json:"errors"`
			}

			switch c.Errors[0].Err.(type) {
			case *validators.ValidationError:
				// Handle validation errors
				err := c.Errors[0].Err.(*validators.ValidationError)
				errorResponse.Status = err.StatusCode
				errorResponse.Errors = err.Errors
			case error:
				// Handle other errors, ensuring a consistent format
				switch c.Errors[0].Type {
				case gin.ErrorTypeBind:
					errorResponse.Status = http.StatusBadRequest
					errorResponse.Errors = []map[string]string{{"error": "Bad Request"}}
				case gin.ErrorTypePrivate:
					// Respect specific status codes set by handlers
					if c.Writer.Status() != 0 {
						errorResponse.Status = c.Writer.Status()
					} else {
						errorResponse.Status = http.StatusInternalServerError
					}
					errorResponse.Errors = []map[string]string{{"error": c.Errors[0].Err.Error()}}
				default:
					errorResponse.Status = http.StatusInternalServerError
					errorResponse.Errors = []map[string]string{{"error": "Internal Server Error"}}
				}
			}

			// Send the structured error response
			c.JSON(errorResponse.Status, errorResponse)
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
