package middlewares

import (
	"encoding/json"
	"net/http"

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

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}

func GinErrorHandlerMiddleware(log *logrus.Logger) gin.HandlerFunc {
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
