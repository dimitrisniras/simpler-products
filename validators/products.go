package validators

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"simpler-products/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ValidateProductID(c *gin.Context) (int, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("invalid product ID")
	}
	if id <= 0 {
		return 0, errors.New("product ID must be positive")
	}
	return id, nil
}

func ValidateProduct(c *gin.Context) (*models.Product, error) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]map[string]string, 0) // Initialize with an empty slice
			for _, fe := range ve {
				errorMsg := getErrorMessage(fe)
				if errorMsg != "" { // Only add if there's a valid error message
					out = append(out, map[string]string{
						"field":   fe.Field(),
						"message": errorMsg,
					})
				}
			}
			return nil, &ValidationError{
				StatusCode: http.StatusBadRequest,
				Errors:     out,
			}
		}
		return nil, err
	}
	return &product, nil
}

// Custom error type for validation errors
type ValidationError struct {
	StatusCode int                 `json:"-"`
	Errors     []map[string]string `json:"errors"`
}

// Helper function to get human-readable error messages
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

// Custom validator for 'gte' tag to handle float64 values
func floatGreaterOrEqualThan(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().Float()
	paramValue := fl.Param()

	// Convert paramValue to float64
	paramFloat, err := strconv.ParseFloat(paramValue, 64)
	if err != nil {
		return false // Handle error if paramValue is not a valid float
	}

	return fieldValue >= paramFloat
}

func init() {
	// Register the custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gte", floatGreaterOrEqualThan)
	}
}

func (v *ValidationError) Error() string {
	return "validation error"
}
