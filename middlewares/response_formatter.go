package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ResponseFormatter(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request and get the data to be sent in the response
		c.Next()

		// Get the data and errors from the context
		data, dataExists := c.Get("data")
		errors, errorsExist := c.Get("errors")

		// Construct the response
		var response struct {
			Status int         `json:"status"`
			Data   interface{} `json:"data,omitempty"`   // Include data only if present
			Errors interface{} `json:"errors,omitempty"` // Include errors only if present
		}

		// Set the status code
		if c.Writer.Status() != 0 {
			response.Status = c.Writer.Status()
		} else {
			response.Status = http.StatusOK // Default to 200 OK if no status is set
		}

		// Include data or errors based on their presence
		if dataExists {
			response.Data = data
		}
		if errorsExist {
			response.Errors = errors
		}

		// Send the formatted response
		c.JSON(response.Status, response)
	}
}
