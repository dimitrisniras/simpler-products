package controllers

import (
	"errors"
	"net/http"
	"simpler-products/models"
	"simpler-products/services"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}

		products, total, err := ps.GetAllProducts(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Set data and pagination in the context
		c.Set("data", products)
		c.Set("pagination", gin.H{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		})
	}
}

func GetProductById(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		product, err := ps.GetProductById(id)
		if err != nil {
			if errors.Is(err, services.ErrProductNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.Error(err)
			}
			return
		}

		// Set data in the context
		c.Set("data", product)
	}
}

func AddProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ps.AddProduct(&product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Set data in the context
		c.Set("data", product)
	}
}

func UpdateProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updatedProduct, err := ps.UpdateProduct(id, &product)
		if err != nil {
			if errors.Is(err, services.ErrProductNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.Error(err)
			}
			return
		}

		// Set data in the context
		c.Set("data", updatedProduct)
	}
}

func PatchProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var product models.Product

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		updatedProduct, err := ps.PatchProduct(id, &product)
		if err != nil {
			if errors.Is(err, services.ErrProductNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.Error(err)
			}
			return
		}

		// Set data in the context
		c.Set("data", updatedProduct)
	}
}

func DeleteProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := ps.DeleteProduct(id); err != nil {
			if errors.Is(err, services.ErrProductNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.Error(err)
			}
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})
	}
}
