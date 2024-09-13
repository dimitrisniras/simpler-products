package validators

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	custom_errors "simpler-products/errors"
	"simpler-products/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateProductID(c *gin.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		c.Status(http.StatusBadRequest)
		c.Set("errors", custom_errors.ErrInvalidProductID)
		return "", custom_errors.ErrInvalidProductID
	}

	return id, nil
}

func ValidateProduct(c *gin.Context) (*models.Product, error) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]map[string]string, 0)
			for _, fe := range ve {
				errorMsg := getErrorMessage(fe)
				if errorMsg != "" {
					out = append(out, map[string]string{
						"message": errorMsg,
					})
				}
			}

			res := &ValidationError{
				Errors: out,
			}

			c.Status(http.StatusBadRequest)
			c.Set("errors", res)
			return nil, res
		}
		c.Set("errors", err)
		return nil, err
	}
	return &product, nil
}

type ValidationError struct {
	Errors []map[string]string `json:"errors"`
}

func (v *ValidationError) Error() string {
	return "validation error"
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

func floatGreaterThan(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().Float()
	paramValue := fl.Param()

	paramFloat, err := strconv.ParseFloat(paramValue, 64)
	if err != nil {
		return false
	}

	return fieldValue > paramFloat
}

func init() {
	validate = validator.New()

	// Register the custom validator
	validate.RegisterValidation("gt", floatGreaterThan)

	// Register custom tag name for 'gt' validator
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name

	})
}
