package controllers

import (
	"errors"
	"net/http"
	"simpler-products/models"
	"simpler-products/services"
	"simpler-products/validators"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := ps.GetAllProducts()
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

func GetProductById(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.ValidateProductID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product, err := ps.GetProductById(id)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func AddProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		product, err := validators.ValidateProduct(c)
		if err != nil {
			c.Error(err)
			return
		}

		if err := ps.AddProduct(product); err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, product)
	}
}

func UpdateProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.ValidateProductID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product, err := validators.ValidateProduct(c)
		if err != nil {
			c.Error(err)
			return
		}

		updatedProduct, err := ps.UpdateProduct(id, product)
		if err != nil {
			if errors.Is(err, services.ErrProductNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.Error(err)
			}
			return
		}

		c.JSON(http.StatusOK, updatedProduct)
	}
}

func PatchProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.ValidateProductID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var product models.Product

		if err := c.ShouldBindJSON(&product); err != nil {
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

		c.JSON(http.StatusOK, updatedProduct)
	}
}

func DeleteProduct(ps services.ProductsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.ValidateProductID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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
