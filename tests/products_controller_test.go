package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"simpler-products/controllers/v1"
	custom_errors "simpler-products/errors"
	"simpler-products/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock implementation of ProductsServiceInterface
type mockProductService struct {
	products []models.Product
	total    int
	err      error
}

func (m *mockProductService) GetAllProducts(limit, offset int) ([]models.Product, int, error) {
	return m.products, m.total, m.err
}

func (m *mockProductService) GetProductById(id string) (*models.Product, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, p := range m.products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, custom_errors.ErrProductNotFound
}

func (m *mockProductService) AddProduct(product *models.Product) error {
	// Simulate ID generation
	product.ID = "generated-uuid"
	m.products = append(m.products, *product)
	m.total++
	return m.err
}

func (m *mockProductService) UpdateProduct(id string, product *models.Product) (*models.Product, error) {
	if m.err != nil {
		return nil, m.err
	}

	for i, p := range m.products {
		if p.ID == id {
			m.products[i] = *product
			return product, nil
		}
	}
	return nil, custom_errors.ErrProductNotFound
}

func (m *mockProductService) DeleteProduct(id string) error {
	if m.err != nil {
		return m.err
	}

	for i, p := range m.products {
		if p.ID == id {
			m.products = append(m.products[:i], m.products[i+1:]...)
			m.total--
			return nil
		}
	}
	return custom_errors.ErrProductNotFound
}

func TestGetAllProductsController(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Create a mock ProductService with sample data
		mockService := &mockProductService{
			products: []models.Product{
				{ID: "uuid1", Name: "Product A", Description: "Description A", Price: 10.99},
				{ID: "uuid2", Name: "Product B", Description: "Description B", Price: 19.95},
			},
			total: 2,
		}

		// Create a request
		req, _ := http.NewRequest("GET", "/products", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.GetAllProducts(mockService)(c)

		data, _ := c.Get("data")
		pagination, _ := c.Get("pagination")
		_, errorsExist := c.Get("errors")

		// Assertions
		assert.Equal(t, errorsExist, false)
		assert.Equal(t, mockService.products, data)
		assert.Equal(t, gin.H(gin.H{"limit": 10, "offset": 0, "total": 2, "count": 2}), pagination)
	})

	t.Run("NoProductsFound", func(t *testing.T) {
		// Create a mock ProductService with an empty product list
		mockService := &mockProductService{
			products: []models.Product{},
			total:    0,
		}

		// Create a request
		req, _ := http.NewRequest("GET", "/products", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.GetAllProducts(mockService)(c)

		data, _ := c.Get("data")
		pagination, _ := c.Get("pagination")
		_, errorsExist := c.Get("errors")

		// Assertions
		assert.Equal(t, errorsExist, false)
		assert.Equal(t, mockService.products, data)
		assert.Equal(t, gin.H(gin.H{"limit": 10, "offset": 0, "total": 0, "count": 0}), pagination)
	})

	t.Run("InvalidLimit", func(t *testing.T) {
		// Create a mock ProductService (not used in this case)
		mockService := &mockProductService{}

		// Create a request with an invalid limit
		req, _ := http.NewRequest("GET", "/products?limit=invalid", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.GetAllProducts(mockService)(c)

		_, dataExists := c.Get("data")
		_, paginationExists := c.Get("pagination")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, paginationExists, false)
		assert.Equal(t, custom_errors.ErrInvalidLimitParameter, err)
	})

	t.Run("InvalidLimitLowerBound", func(t *testing.T) {
		// Create a mock ProductService (not used in this case)
		mockService := &mockProductService{}

		// Create a request with an invalid limit
		req, _ := http.NewRequest("GET", "/products?limit=-1", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.GetAllProducts(mockService)(c)

		_, dataExists := c.Get("data")
		_, paginationExists := c.Get("pagination")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, paginationExists, false)
		assert.Equal(t, custom_errors.ErrInvalidLimitParameter, err)
	})

	t.Run("InvalidLimitLowerBound", func(t *testing.T) {
		// Create a mock ProductService (not used in this case)
		mockService := &mockProductService{}

		// Create a request with an invalid limit
		req, _ := http.NewRequest("GET", "/products?limit=101", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.GetAllProducts(mockService)(c)

		_, dataExists := c.Get("data")
		_, paginationExists := c.Get("pagination")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, paginationExists, false)
		assert.Equal(t, custom_errors.ErrInvalidLimitParameter, err)
	})

	t.Run("InvalidOffset", func(t *testing.T) {
		// Create a mock ProductService (not used in this case)
		mockService := &mockProductService{}

		// Create a request with an invalid offset
		req, _ := http.NewRequest("GET", "/products?offset=-1", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.GetAllProducts(mockService)(c)

		_, dataExists := c.Get("data")
		_, paginationExists := c.Get("pagination")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, paginationExists, false)
		assert.Equal(t, custom_errors.ErrInvalidOffsetParameter, err)
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Create a mock ProductService that returns an error
		mockError := errors.New("service error")
		mockService := &mockProductService{
			err: mockError,
		}

		// Create a request
		req, _ := http.NewRequest("GET", "/products", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.GetAllProducts(mockService)(c)

		_, dataExists := c.Get("data")
		_, paginationExists := c.Get("pagination")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, dataExists, false)
		assert.Equal(t, paginationExists, false)
		assert.Equal(t, mockError, err)
	})
}

func TestGetProductByIdController(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Create a mock ProductService with sample data
		mockService := &mockProductService{
			products: []models.Product{
				{ID: "uuid1", Name: "Product A", Description: "Description A", Price: 10.99},
				{ID: "uuid2", Name: "Product B", Description: "Description B", Price: 19.95},
			},
		}

		// Create a request
		req, _ := http.NewRequest("GET", "/products/uuid1", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "uuid1"}}

		// Call the handler function
		controllers.GetProductById(mockService)(c)

		data, _ := c.Get("data")
		_, errorsExist := c.Get("errors")

		// Assertions
		assert.Equal(t, errorsExist, false)
		assert.Equal(t, [1]*models.Product{&mockService.products[0]}, data)
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		// Create a mock ProductService with sample data
		mockService := &mockProductService{
			products: []models.Product{
				{ID: "uuid1", Name: "Product A", Description: "Description A", Price: 10.99},
				{ID: "uuid2", Name: "Product B", Description: "Description B", Price: 19.95},
			},
		}

		// Create a request with a non-existent ID
		req, _ := http.NewRequest("GET", "/products/non_existent_id", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "non_existent_id"}}

		// Call the handler function
		controllers.GetProductById(mockService)(c)

		_, dataExists := c.Get("data")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusNotFound, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, custom_errors.ErrProductNotFound, err)
	})

	t.Run("InvalidID_Empty", func(t *testing.T) {
		// Create a mock ProductService (not used in this case)
		mockService := &mockProductService{}

		// Create a request with an empty ID
		req, _ := http.NewRequest("GET", "/products/", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

		// Call the handler function
		controllers.GetProductById(mockService)(c)

		_, dataExists := c.Get("data")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, custom_errors.ErrInvalidProductID, err)
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Create a mock ProductService that returns an error
		mockError := errors.New("service error")
		mockService := &mockProductService{
			err: mockError,
		}

		// Create a request
		req, _ := http.NewRequest("GET", "/products/uuid1", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "uuid1"}}

		// Call the handler function
		controllers.GetProductById(mockService)(c)

		_, dataExists := c.Get("data")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, dataExists, false)
		assert.Equal(t, mockError, err)
	})
}

func TestAddProductController(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Create a mock ProductService with an empty product list
		mockService := &mockProductService{
			products: []models.Product{},
		}

		// Create a request with product data
		productData := &models.Product{
			Name:        "New Product",
			Description: "Description",
			Price:       9.99,
		}

		jsonData, _ := json.Marshal(productData)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.AddProduct(mockService)(c)

		data, _ := c.Get("data")
		_, errorsExist := c.Get("errors")

		fmt.Println(data)

		// Assertions
		assert.Equal(t, http.StatusCreated, c.Writer.Status())
		assert.Equal(t, errorsExist, false)
		// assert.Equal(t, productData, data)
	})

	t.Run("InvalidProductData", func(t *testing.T) {
		// Create a mock ProductService (not used in this case)
		mockService := &mockProductService{}

		// Create a request with invalid product data
		invalidProductData := map[string]interface{}{
			"name":  "",   // Missing required field
			"price": -5.0, // Invalid price
		}
		jsonData, _ := json.Marshal(invalidProductData)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.AddProduct(mockService)(c)

		_, dataExists := c.Get("data")
		_, errExist := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, errExist, true)
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Create a mock ProductService that returns an error
		mockError := errors.New("service error")
		mockService := &mockProductService{
			err: mockError,
		}

		// Create a request with product data
		productData := &models.Product{
			Name:        "New Product",
			Description: "Description",
			Price:       9.99,
		}
		jsonData, _ := json.Marshal(productData)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		controllers.AddProduct(mockService)(c)

		_, dataExists := c.Get("data")
		err, errExist := c.Get("errors")

		// Assertions
		assert.Equal(t, dataExists, false)
		assert.Equal(t, errExist, true)
		assert.Equal(t, mockError, err)
	})
}

func TestUpdateProductController(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Create a mock ProductService with sample data
		mockService := &mockProductService{
			products: []models.Product{
				{ID: "uuid1", Name: "Product A", Description: "Description A", Price: 10.99},
				{ID: "uuid2", Name: "Product B", Description: "Description B", Price: 19.95},
			},
		}

		// Create a request with updated product data
		updatedProduct := &models.Product{
			Name:        "Updated Product A",
			Description: "Updated Description A",
			Price:       12.99,
		}
		jsonData, _ := json.Marshal(updatedProduct)
		req, _ := http.NewRequest("PUT", "/products/uuid1", bytes.NewBuffer(jsonData))

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "uuid1"}}

		// Call the handler function
		controllers.UpdateProduct(mockService)(c)

		data, _ := c.Get("data")
		_, errorsExist := c.Get("errors")

		// Assertions
		assert.Equal(t, errorsExist, false)
		assert.Equal(t, [1]*models.Product{updatedProduct}, data)
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		// Create a mock ProductService
		mockService := &mockProductService{}

		// Create a request with a non-existent ID
		updatedProduct := &models.Product{
			Name:        "Updated Product A",
			Description: "Updated Description A",
			Price:       12.99,
		}
		jsonData, _ := json.Marshal(updatedProduct)
		req, _ := http.NewRequest("PUT", "/products/non_existent_id", bytes.NewBuffer(jsonData))

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "non_existent_id"}}

		// Call the handler function
		controllers.UpdateProduct(mockService)(c)

		_, dataExists := c.Get("data")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusNotFound, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, custom_errors.ErrProductNotFound, err)
	})

	t.Run("InvalidID_Empty", func(t *testing.T) {
		// Create a mock ProductService (not used in this case)
		mockService := &mockProductService{}

		// Create a request with an empty ID
		updatedProduct := &models.Product{
			Name:        "Updated Product A",
			Description: "Updated Description A",
			Price:       12.99,
		}
		jsonData, _ := json.Marshal(updatedProduct)
		req, _ := http.NewRequest("PUT", "/products/", bytes.NewBuffer(jsonData))

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

		// Call the handler function
		controllers.UpdateProduct(mockService)(c)

		_, dataExists := c.Get("data")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, custom_errors.ErrInvalidProductID, err)
	})

	t.Run("InvalidProductData", func(t *testing.T) {
		// Create a mock ProductService (not used in this case)
		mockService := &mockProductService{}

		// Create a request with invalid product data
		invalidProductData := map[string]interface{}{
			"name":  "",   // Missing required field
			"price": -5.0, // Invalid price
		}
		jsonData, _ := json.Marshal(invalidProductData)
		req, _ := http.NewRequest("PUT", "/products/uuid1", bytes.NewBuffer(jsonData))

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "uuid1"}}

		// Call the handler function
		controllers.UpdateProduct(mockService)(c)

		_, dataExists := c.Get("data")
		_, errExists := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, errExists, true)
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Create a mock ProductService that returns an error
		mockError := errors.New("service error")
		mockService := &mockProductService{
			err: mockError,
		}

		// Create a request with updated product data
		updatedProduct := &models.Product{
			Name:        "Updated Product A",
			Description: "Updated Description A",
			Price:       12.99,
		}
		jsonData, _ := json.Marshal(updatedProduct)
		req, _ := http.NewRequest("PUT", "/products/uuid1", bytes.NewBuffer(jsonData))

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "uuid1"}}

		// Call the handler function
		controllers.UpdateProduct(mockService)(c)

		_, dataExists := c.Get("data")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, dataExists, false)
		assert.Equal(t, mockError, err)
	})
}

func TestDeleteProductController(t *testing.T) {
	// Set Gin to TestMode
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Create a mock ProductService with sample data
		mockService := &mockProductService{
			products: []models.Product{
				{ID: "uuid1", Name: "Product A", Description: "Description A", Price: 10.99},
				{ID: "uuid2", Name: "Product B", Description: "Description B", Price: 19.95},
			},
		}

		// Create a request
		req, _ := http.NewRequest("DELETE", "/products/uuid1", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "uuid1"}}

		// Call the handler function
		controllers.DeleteProduct(mockService)(c)

		_, dataExists := c.Get("data")
		_, errExists := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusNoContent, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, errExists, false)
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		// Create a mock ProductService
		mockService := &mockProductService{}

		// Create a request with a non-existent ID
		req, _ := http.NewRequest("DELETE", "/products/non_existent_id", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "non_existent_id"}}

		// Call the handler function
		controllers.DeleteProduct(mockService)(c)

		_, dataExists := c.Get("data")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusNotFound, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, custom_errors.ErrProductNotFound, err)
	})

	t.Run("InvalidID_Empty", func(t *testing.T) {
		// Create a mock ProductService (not used in this case)
		mockService := &mockProductService{}

		// Create a request with an empty ID
		req, _ := http.NewRequest("DELETE", "/products/", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

		// Call the handler function
		controllers.DeleteProduct(mockService)(c)

		_, dataExists := c.Get("data")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
		assert.Equal(t, dataExists, false)
		assert.Equal(t, custom_errors.ErrInvalidProductID, err)
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Create a mock ProductService that returns an error
		mockError := errors.New("service error")
		mockService := &mockProductService{
			err: mockError,
		}

		// Create a request
		req, _ := http.NewRequest("DELETE", "/products/uuid1", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Create a Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "uuid1"}}

		// Call the handler function
		controllers.DeleteProduct(mockService)(c)

		_, dataExists := c.Get("data")
		err, _ := c.Get("errors")

		// Assertions
		assert.Equal(t, dataExists, false)
		assert.Equal(t, mockError, err)
	})
}
