package tests

import (
	"database/sql"
	"errors"
	"simpler-products/models"
	"simpler-products/services"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	custom_errors "simpler-products/errors"
)

func TestGetAllProductsService(t *testing.T) {
	// Set up mock database
	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a mock ProductsService
	log := logrus.New()
	productService := &services.ProductsService{
		DB:  db,
		Log: log,
	}

	t.Run("Success", func(t *testing.T) {
		// Mock the count query
		dbMock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Products").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// Mock the paginated query
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
			AddRow("uuid1", "Product A", "Description A", 10.99).
			AddRow("uuid2", "Product B", "Description B", 19.95)

		dbMock.ExpectQuery("SELECT \\* FROM Products LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		// Call the service function
		products, total, err := productService.GetAllProducts(10, 0)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, 2, len(products))
		assert.Equal(t, 2, total)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DifferentLimitAndOffsetCombinations", func(t *testing.T) {
		testCases := []struct {
			name   string
			limit  int
			offset int
			total  int
			rows   *sqlmock.Rows
		}{
			{
				name:   "FirstPage",
				limit:  5,
				offset: 0,
				total:  10,
				rows: sqlmock.NewRows([]string{"id", "name", "description", "price"}).
					AddRow("uuid1", "Product A", "Description A", 10.99).
					AddRow("uuid2", "Product B", "Description B", 19.95).
					AddRow("uuid3", "Product C", "Description C", 5.50).
					AddRow("uuid4", "Product D", "Description D", 8.25).
					AddRow("uuid5", "Product E", "Description E", 15.00),
			},
			{
				name:   "SecondPage",
				limit:  5,
				offset: 5,
				total:  10,
				rows: sqlmock.NewRows([]string{"id", "name", "description", "price"}).
					AddRow("uuid6", "Product F", "Description F", 7.75).
					AddRow("uuid7", "Product G", "Description G", 22.30).
					AddRow("uuid8", "Product H", "Description H", 3.15).
					AddRow("uuid9", "Product I", "Description I", 11.80).
					AddRow("uuid10", "Product J", "Description J", 6.40),
			},
			{
				name:   "PartialLastPage",
				limit:  5,
				offset: 10,
				total:  12,
				rows: sqlmock.NewRows([]string{"id", "name", "description", "price"}).
					AddRow("uuid11", "Product K", "Description K", 9.00).
					AddRow("uuid12", "Product L", "Description L", 4.60),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Mock the count query
				dbMock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Products").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tc.total))

				// Mock the paginated query
				dbMock.ExpectQuery("SELECT \\* FROM Products LIMIT \\? OFFSET \\?").
					WithArgs(tc.limit, tc.offset).
					WillReturnRows(tc.rows)

				// Call the service function
				products, total, err := productService.GetAllProducts(tc.limit, tc.offset)

				// Assertions
				assert.NoError(t, err)
				assert.Equal(t, len(products), len(products))
				assert.Equal(t, tc.total, total)

				// Ensure all expectations were met
				if err := dbMock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			})
		}
	})

	t.Run("LargeOffsetExceedingTotal", func(t *testing.T) {
		// Mock a successful count query
		dbMock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Products").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		// Mock an empty result set for a large offset
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"})
		dbMock.ExpectQuery("SELECT \\* FROM Products LIMIT \\? OFFSET \\?").
			WithArgs(10, 100). // Large offset
			WillReturnRows(rows)

		// Call the service function
		products, total, err := productService.GetAllProducts(10, 100)

		// Assertions
		assert.NoError(t, err)
		assert.Empty(t, products) // No products should be returned
		assert.Equal(t, 5, total) // Total count should still be accurate

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("ErrorFetchingProducts", func(t *testing.T) {
		// Mock a successful count query
		dbMock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Products").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		// Mock an error when fetching the paginated products
		dbMock.ExpectQuery("SELECT \\* FROM Products LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnError(errors.New("database error"))

		// Call the service function
		products, total, err := productService.GetAllProducts(10, 0)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, total) // Total should be 0 if there's an error fetching products

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("NoProductsFound", func(t *testing.T) {
		// Mock the count query
		dbMock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Products").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Mock an empty result set
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"})
		dbMock.ExpectQuery("SELECT (.+) FROM Products").WillReturnRows(rows)

		// Call the service function
		products, total, err := productService.GetAllProducts(10, 0)

		// Assertions
		assert.NoError(t, err)
		assert.Empty(t, products)
		assert.Equal(t, 0, total)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("ErrorScanningRow", func(t *testing.T) {
		// Mock a successful count query
		dbMock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Products").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// Mock the paginated query to return rows with an incompatible data type
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
			AddRow("uuid1", "Product A", "Description A", "invalid_price") // Invalid price format

		dbMock.ExpectQuery("SELECT \\* FROM Products LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		// Call the service function
		products, total, err := productService.GetAllProducts(10, 0)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, total) // Total should be 0 if there's an error scanning rows

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("ErrorOnCountQuery", func(t *testing.T) {
		// Mock an error during the count query
		dbMock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Products").
			WillReturnError(errors.New("database error during count"))

		// Call the service function
		products, total, err := productService.GetAllProducts(10, 0)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, total)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("InvalidLimit", func(t *testing.T) {
		// Call the service function with an invalid limit
		products, total, err := productService.GetAllProducts(0, 0)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, total)
	})

	t.Run("InvalidOffset", func(t *testing.T) {
		// Call the service function with an invalid offset
		products, total, err := productService.GetAllProducts(10, -1)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, total)
	})
}

func TestGetProductByIdService(t *testing.T) {
	// Set up mock database
	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a mock ProductsService
	log := logrus.New()
	productService := &services.ProductsService{
		DB:  db,
		Log: log,
	}

	t.Run("Success", func(t *testing.T) {
		// Mock the database query to return a product
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
			AddRow("uuid1", "Product A", "Description A", 10.99)

		dbMock.ExpectQuery("SELECT (.+) FROM Products WHERE id = ?").
			WithArgs("uuid1").
			WillReturnRows(rows)

		// Call the service function
		product, err := productService.GetProductById("uuid1")

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "uuid1", product.ID)
		assert.Equal(t, "Product A", product.Name)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		// Mock the database query to return no rows
		dbMock.ExpectQuery("SELECT (.+) FROM Products WHERE id = ?").
			WithArgs("non_existent_id").
			WillReturnError(sql.ErrNoRows)

		// Call the service function
		product, err := productService.GetProductById("non_existent_id")

		// Assertions
		assert.Error(t, err)
		assert.True(t, errors.Is(err, custom_errors.ErrProductNotFound))
		assert.Nil(t, product)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("InvalidIdFormat", func(t *testing.T) {
		// Call the service function with an id that has an invalid format
		product, err := productService.GetProductById("invalid_id")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, product)

		// Ensure all expectations were met (no database interactions should occur in this case)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Mock a database error
		dbMock.ExpectQuery("SELECT (.+) FROM Products WHERE id = ?").
			WithArgs("uuid1").
			WillReturnError(errors.New("database error"))

		// Call the service function
		product, err := productService.GetProductById("uuid1")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, product)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseErrorScanningRow", func(t *testing.T) {
		// Mock the database query to return a row, but simulate an error during scanning
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
			AddRow("uuid1", "Product A", "Description A", "invalid_price") // Invalid price format

		dbMock.ExpectQuery("SELECT (.+) FROM Products WHERE id = ?").
			WithArgs("uuid1").
			WillReturnRows(rows)

		// Call the service function
		product, err := productService.GetProductById("uuid1")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, product)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("VeryLongId", func(t *testing.T) {
		// Create a very long id
		veryLongID := strings.Repeat("a", 1000)

		// Mock the database query with the very long ID
		dbMock.ExpectQuery("SELECT (.+) FROM Products WHERE id = ?").
			WithArgs(veryLongID).
			WillReturnError(sql.ErrNoRows)

		// Call the service function
		product, err := productService.GetProductById(veryLongID)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, product)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestAddProductService(t *testing.T) {
	// Set up mock database
	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a mock ProductsService
	log := logrus.New()
	productService := &services.ProductsService{
		DB:  db,
		Log: log,
	}

	t.Run("Success", func(t *testing.T) {
		// Mock the database Exec to return a successful result
		dbMock.ExpectExec("INSERT INTO Products \\(id, name, description, price\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
			WithArgs(sqlmock.AnyArg(), "New Product", "Description", 9.99). // sqlmock.AnyArg() for the UUID
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Create a new product (without an ID, as it will be generated)
		newProduct := &models.Product{
			Name:        "New Product",
			Description: "Description",
			Price:       9.99,
		}

		// Call the service function
		err := productService.AddProduct(newProduct)

		// Assertions
		assert.NoError(t, err)
		assert.NotEmpty(t, newProduct.ID) // Check if the ID is set after creation (UUID generated)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Mock a database error during insertion
		dbMock.ExpectExec("INSERT INTO Products \\(id, name, description, price\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
			WithArgs(sqlmock.AnyArg(), "New Product", "Description", 9.99).
			WillReturnError(errors.New("database error"))

		// Create a new product
		newProduct := &models.Product{
			Name:        "New Product",
			Description: "Description",
			Price:       9.99,
		}

		// Call the service function
		err := productService.AddProduct(newProduct)

		// Assertions
		assert.Error(t, err)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("InvalidProductData", func(t *testing.T) {
		// Create a new product with invalid data
		invalidProduct := &models.Product{
			// Missing Name and Description
			Price: -5.0, // Invalid price
		}

		// Call the service function
		err := productService.AddProduct(invalidProduct)

		// Assertions
		assert.Error(t, err)
	})

	t.Run("DuplicateProductId", func(t *testing.T) {
		// Mock the database Exec to simulate a duplicate key error
		dbMock.ExpectExec("INSERT INTO Products \\(id, name, description, price\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
			WithArgs(sqlmock.AnyArg(), "New Product", "Description", 9.99).
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'some-uuid' for key 'PRIMARY'"})

		// Create a new product
		newProduct := &models.Product{
			Name:        "New Product",
			Description: "Description",
			Price:       9.99,
		}

		// Call the service function
		err := productService.AddProduct(newProduct)

		// Assertions
		assert.Error(t, err)
	})

	t.Run("DifferentValidProductData", func(t *testing.T) {
		testCases := []struct {
			name        string
			productData *models.Product
		}{
			{
				name: "BasicProduct",
				productData: &models.Product{
					Name:        "Basic Product",
					Description: "Simple description",
					Price:       10.0,
				},
			},
			{
				name: "NormalProduct",
				productData: &models.Product{
					Name:        "Normal Item",
					Description: "Normal product with normal price",
					Price:       100.0,
				},
			},
			{
				name: "ExpensiveProduct",
				productData: &models.Product{
					Name:        "Luxury Item",
					Description: "High-end product with premium features",
					Price:       1000.0,
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Mock the database Exec to return a successful result
				dbMock.ExpectExec("INSERT INTO Products \\(id, name, description, price\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
					WithArgs(sqlmock.AnyArg(), tc.productData.Name, tc.productData.Description, tc.productData.Price).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Call the service function
				err := productService.AddProduct(tc.productData)

				// Assertions
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.productData.ID)

				// Ensure all expectations were met
				if err := dbMock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			})
		}
	})

	t.Run("VeryLongNameOrDescription", func(t *testing.T) {
		veryLongName := strings.Repeat("a", 256)
		veryLongDescription := strings.Repeat("b", 1024)

		newProduct := &models.Product{
			Name:        veryLongName,
			Description: veryLongDescription,
			Price:       9.99,
		}

		// Call the service function
		err := productService.AddProduct(newProduct)

		assert.Error(t, err)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DifferentPriceValues", func(t *testing.T) {
		testCases := []struct {
			name  string
			price float64
			valid bool
		}{
			{
				name:  "ZeroPrice",
				price: 0.0,
				valid: false,
			},
			{
				name:  "VeryLargePrice",
				price: 1e10,
				valid: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Mock the database Exec only if the price is valid
				if tc.valid {
					dbMock.ExpectExec("INSERT INTO Products \\(id, name, description, price\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
						WithArgs(sqlmock.AnyArg(), "Product", "Description", tc.price).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				// Create a new product
				newProduct := &models.Product{
					Name:        "Product",
					Description: "Description",
					Price:       tc.price,
				}

				// Call the service function
				err := productService.AddProduct(newProduct)

				// Assertions
				if tc.valid {
					assert.NoError(t, err)
					assert.NotEmpty(t, newProduct.ID)
				} else {
					assert.Error(t, err)
				}

				// Ensure all expectations were met
				if err := dbMock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			})
		}
	})
}

func TestUpdateProductService(t *testing.T) {
	// Set up mock database
	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a mock ProductsService
	log := logrus.New()
	productService := &services.ProductsService{
		DB:  db,
		Log: log,
	}

	t.Run("Success", func(t *testing.T) {
		// Mock the database Exec for update and the query to fetch the updated product
		dbMock.ExpectExec("UPDATE Products SET name = \\?, description = \\?, price = \\? WHERE id = \\?").
			WithArgs("Updated Product", "Updated Description", 12.99, "uuid1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
			AddRow("uuid1", "Updated Product", "Updated Description", 12.99)

		dbMock.ExpectQuery("SELECT (.+) FROM Products WHERE id = ?").
			WithArgs("uuid1").
			WillReturnRows(rows)

		// Create an updated product
		updatedProduct := &models.Product{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       12.99,
		}

		// Call the service function
		product, err := productService.UpdateProduct("uuid1", updatedProduct)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "uuid1", product.ID)
		assert.Equal(t, "Updated Product", product.Name)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		// Mock the database Exec to return no rows affected (product not found)
		dbMock.ExpectExec("UPDATE Products SET name = \\?, description = \\?, price = \\? WHERE id = \\?").
			WithArgs("Updated Product", "Updated Description", 12.99, "non_existent_id").
			WillReturnResult(sqlmock.NewResult(0, 0))

		dbMock.ExpectQuery("SELECT (.+) FROM Products WHERE id = ?").
			WithArgs("non_existent_id").
			WillReturnError(sql.ErrNoRows)

		// Create an updated product
		updatedProduct := &models.Product{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       12.99,
		}

		// Call the service function
		product, err := productService.UpdateProduct("non_existent_id", updatedProduct)

		// Assertions
		assert.Error(t, err)
		assert.True(t, errors.Is(err, custom_errors.ErrProductNotFound))
		assert.Nil(t, product)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Mock a database error during the update
		dbMock.ExpectExec("UPDATE Products SET name = \\?, description = \\?, price = \\? WHERE id = \\?").
			WithArgs("Updated Product", "Updated Description", 12.99, "uuid1").
			WillReturnError(errors.New("database error"))

		// Create an updated product
		updatedProduct := &models.Product{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       12.99,
		}

		// Call the service function
		product, err := productService.UpdateProduct("uuid1", updatedProduct)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, product)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("InvalidProductData", func(t *testing.T) {
		// Create an updated product with invalid data
		invalidProduct := &models.Product{
			// Missing Name and Description
			Price: -5.0, // Invalid price
		}

		// Call the service function
		product, err := productService.UpdateProduct("uuid1", invalidProduct)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, product)

		// Ensure all expectations were met (no database interactions should occur in this case)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseErrorFetchingUpdatedProduct", func(t *testing.T) {
		// Mock a successful update
		dbMock.ExpectExec("UPDATE Products SET name = \\?, description = \\?, price = \\? WHERE id = \\?").
			WithArgs("Updated Product", "Updated Description", 12.99, "uuid1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock an error when fetching the updated product
		dbMock.ExpectQuery("SELECT \\* FROM Products WHERE id = \\?").
			WithArgs("uuid1").
			WillReturnError(errors.New("database error fetching updated product"))

		// Create an updated product
		updatedProduct := &models.Product{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       12.99,
		}

		// Call the service function
		product, err := productService.UpdateProduct("uuid1", updatedProduct)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, product)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestDeleteProductService(t *testing.T) {
	// Set up mock database
	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a mock ProductsService
	log := logrus.New()
	productService := &services.ProductsService{
		DB:  db,
		Log: log,
	}

	t.Run("Success", func(t *testing.T) {
		// Mock the database query to ensure the product exists before deletion
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
			AddRow("uuid1", "Product A", "Description A", 10.99)

		dbMock.ExpectQuery("SELECT \\* FROM Products WHERE id = \\?").
			WithArgs("uuid1").
			WillReturnRows(rows)

		// Mock the database Exec for deletion
		dbMock.ExpectExec("DELETE FROM Products WHERE id = \\?").
			WithArgs("uuid1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Call the service function
		err := productService.DeleteProduct("uuid1")

		// Assertions
		assert.NoError(t, err)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		// Mock the database query to return no rows (product not found)
		dbMock.ExpectQuery("SELECT \\* FROM Products WHERE id = \\?").
			WithArgs("non_existent_id").
			WillReturnError(sql.ErrNoRows)

		// Call the service function
		err := productService.DeleteProduct("non_existent_id")

		// Assertions
		assert.Error(t, err)
		assert.True(t, errors.Is(err, custom_errors.ErrProductNotFound))

		// Ensure all expectations were met (no DELETE should be executed)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseErrorDuringFetch", func(t *testing.T) {
		// Mock a database error during the product fetch
		dbMock.ExpectQuery("SELECT \\* FROM Products WHERE id = \\?").
			WithArgs("uuid1").
			WillReturnError(errors.New("database error"))

		// Call the service function
		err := productService.DeleteProduct("uuid1")

		// Assertions
		assert.Error(t, err)

		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseErrorDuringDelete", func(t *testing.T) {
		// Mock a successful product fetch
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
			AddRow("uuid1", "Product A", "Description A", 10.99)

		dbMock.ExpectQuery("SELECT \\* FROM Products WHERE id = \\?").
			WithArgs("uuid1").
			WillReturnRows(rows)

		// Mock a database error during the delete operation
		dbMock.ExpectExec("DELETE FROM Products WHERE id = \\?").
			WithArgs("uuid1").
			WillReturnError(errors.New("database error"))

		// Call the service function
		err := productService.DeleteProduct("uuid1")

		// Assertions
		assert.Error(t, err)

		// Ensure all expectations were met
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
