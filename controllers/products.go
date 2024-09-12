package controllers

import (
	"net/http"
	"simpler-products/models"
	"simpler-products/services"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	}
}

func GetProductById(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
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

		c.JSON(http.StatusCreated, product)
	}
}

func UpdateProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	}
}

func PatchProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	}
}

func DeleteProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	}
}
