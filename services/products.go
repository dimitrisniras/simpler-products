package services

import (
	"database/sql"
	custom_errors "simpler-products/errors"
	"simpler-products/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductsServiceInterface interface {
	GetAllProducts(limit, offset int) ([]models.Product, int, error)
	GetProductById(id string) (*models.Product, error)
	AddProduct(product *models.Product) error
	UpdateProduct(id string, product *models.Product) (*models.Product, error)
	DeleteProduct(id string) error
}

type ProductsService struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func (ps *ProductsService) GetAllProducts(limit, offset int) ([]models.Product, int, error) {
	ps.Log.Debugf("Fetching products from database, limit: %d, offset: %d", limit, offset)

	// 1. Get the total count of products
	var totalCount int
	err := ps.DB.QueryRow("SELECT COUNT(*) FROM Products").Scan(&totalCount)
	if err != nil {
		ps.Log.Errorf("Error getting total product count: %v", err)
		return nil, 0, err
	}

	// 2. Fetch paginated products
	rows, err := ps.DB.Query("SELECT * FROM Products LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		ps.Log.Errorf("Error fetching products: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			ps.Log.Errorf("Error scanning product row: %v", err)
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, totalCount, nil
}

func (ps *ProductsService) GetProductById(id string) (*models.Product, error) {
	ps.Log.Debugf("Fetching product with ID: %v from database", id)

	var product models.Product
	err := ps.DB.QueryRow("SELECT * FROM Products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrProductNotFound
		}
		ps.Log.Errorf("Error fetching product: %v", err)
		return nil, err
	}

	return &product, nil
}

func (ps *ProductsService) AddProduct(product *models.Product) error {
	ps.Log.Debugf("Creating new product in database, data: %+v", product)

	uuid := uuid.NewString()
	_, err := ps.DB.Exec("INSERT INTO Products (id, name, description, price) VALUES (?, ?, ?, ?)", uuid, product.Name, product.Description, product.Price)
	if err != nil {
		ps.Log.Errorf("Error creating new product: %v", err)
		return err
	}

	product.ID = uuid

	return nil
}

func (ps *ProductsService) UpdateProduct(id string, product *models.Product) (*models.Product, error) {
	ps.Log.Debugf("Updating product with ID: %v in database, data: %+v", id, product)

	_, err := ps.DB.Exec("UPDATE Products SET name = ?, description = ?, price = ? WHERE id = ?", product.Name, product.Description, product.Price, id)
	if err != nil {
		ps.Log.Errorf("Error updating product: %v", err)
		return nil, err
	}

	// Fetch the updated product from the database
	updatedProduct, err := ps.GetProductById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrProductNotFound
		}
		ps.Log.Errorf("Error updating product: %v", err)
		return nil, err
	}

	return updatedProduct, nil
}

func (ps *ProductsService) DeleteProduct(id string) error {
	ps.Log.Debugf("Deleting product with ID: %v from database", id)

	// Fetch the product to be deleted from the database
	_, err := ps.GetProductById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return custom_errors.ErrProductNotFound
		}
		ps.Log.Errorf("Error deleting product: %v", err)
		return err
	}

	_, err = ps.DB.Exec("DELETE FROM Products WHERE id = ?", id)
	if err != nil {
		ps.Log.Errorf("Error deleting product: %v", err)
		return err
	}

	return nil
}
