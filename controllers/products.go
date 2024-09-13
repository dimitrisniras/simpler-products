package controllers

import (
	"errors"
	"net/http"
	"simpler-products/models"
	"simpler-products/services"
	"simpler-products/validators"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get pagination parameters from query string
		limitStr := c.DefaultQuery("limit", "10")  // Default limit is 10
		offsetStr := c.DefaultQuery("offset", "0") // Default offset is 0

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.Status(http.StatusBadRequest)
			c.Set("errors", errors.New("invalid limit parameter"))
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.Status(http.StatusBadRequest)
			c.Set("errors", errors.New("invalid offset parameter"))
			return
		}

		products, total, err := ps.GetAllProducts(limit, offset)
		if err != nil {
			c.Set("errors", err)
			return
		}

		// Set data and pagination in the context
		c.Set("data", products)
		c.Set("pagination", gin.H{
			"limit":  limit,
			"offset": offset,
			"total":  total,
			"count":  len(products),
		})
	}
}

func GetProductById(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.ValidateProductID(c)
		if err != nil {
			return
		}

		product, err := ps.GetProductById(id)
		if err != nil {
			if errors.Is(err, services.ErrProductNotFound) {
				c.Status(http.StatusNotFound)
			}
			c.Set("errors", err)
			return
		}

		// Set data in the context
		c.Set("data", [1]*models.Product{product})
	}
}

func AddProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		product, err := validators.ValidateProduct(c)
		if err != nil {
			return
		}

		if err := ps.AddProduct(product); err != nil {
			c.Set("errors", err)
			return
		}

		// Set data in the context
		c.Status(http.StatusCreated)
		c.Set("data", [1]*models.Product{product})
	}
}

func UpdateProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.ValidateProductID(c)
		if err != nil {
			return
		}

		product, err := validators.ValidateProduct(c)
		if err != nil {
			return
		}

		updatedProduct, err := ps.UpdateProduct(id, product)
		if err != nil {
			if errors.Is(err, services.ErrProductNotFound) {
				c.Status(http.StatusNotFound)
			}
			c.Set("errors", err)
			return
		}

		// Set data in the context
		c.Set("data", [1]*models.Product{updatedProduct})
	}
}

func DeleteProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.ValidateProductID(c)
		if err != nil {
			return
		}

		if err := ps.DeleteProduct(id); err != nil {
			if errors.Is(err, services.ErrProductNotFound) {
				c.Status(http.StatusNotFound)
			}
			c.Set("errors", err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
