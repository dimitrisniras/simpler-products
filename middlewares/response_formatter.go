package middlewares

import (
	"net/http"
	"simpler-products/validators"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ResponseFormatter(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request and get the data to be sent in the response
		c.Next()

		// Get the data and errors from the context
		data, dataExists := c.Get("data")
		pagination, paginationExists := c.Get("pagination")
		errors, errorsExist := c.Get("errors")

		// Construct the response
		var response struct {
			Status     int         `json:"status"`
			Data       interface{} `json:"data,omitempty"`       // Include data only if present
			Pagination interface{} `json:"pagination,omitempty"` // Include pagination only if present
			Errors     interface{} `json:"errors,omitempty"`     // Include errors only if present
		}

		// Set the status code based on the presence of errors
		if errorsExist {
			if c.Writer.Status() != 0 {
				response.Status = c.Writer.Status()
			} else {
				response.Status = http.StatusInternalServerError // Default to 500 if no status is set
			}

			// Log the error
			log.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())

			// Format the errors consistently
			formattedErrors := make([]map[string]any, 0)
			switch errors.(type) {
			case *validators.ValidationError:
				err := errors.(*validators.ValidationError)
				for _, validationError := range err.Errors {
					for _, errorMsg := range validationError {
						formattedErrors = append(formattedErrors, map[string]any{
							"message": errorMsg,
						})
					}
				}
			case []string:
				for _, errorMsg := range errors.([]string) {
					formattedErrors = append(formattedErrors, map[string]any{
						"message": errorMsg,
					})
				}
			default:
				formattedErrors = append(formattedErrors, map[string]any{
					"message": errors.(error).Error(),
				})
			}
			response.Errors = formattedErrors
		} else {
			if c.Writer.Status() != 0 {
				response.Status = c.Writer.Status()
			} else {
				response.Status = http.StatusOK // Default to 200 OK if no status is set
			}

			// Include data only if present and the status code is 200 OK
			if dataExists {
				response.Data = data
			}

			// Include pagination only if present and the status code is 200 OK
			if paginationExists {
				response.Pagination = pagination
			}
		}

		// Send the formatted response
		c.JSON(response.Status, response)
	}
}
