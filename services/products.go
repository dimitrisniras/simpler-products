package services

import (
	"database/sql"
	"errors"
	"fmt"
	"simpler-products/models"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductsServiceInterface interface {
	GetAllProducts() ([]models.Product, error)
	GetProductById(id int) (*models.Product, error)
	AddProduct(product *models.Product) error
	UpdateProduct(id int, product *models.Product) (*models.Product, error)
	PatchProduct(id int, product *models.Product) (*models.Product, error)
	DeleteProduct(id int) error
}

type ProductsService struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func (ps *ProductsService) GetAllProducts() ([]models.Product, error) {
	ps.Log.Debug("Fetching all products from database")

	rows, err := ps.DB.Query("SELECT * FROM Products")
	if err != nil {
		ps.Log.Errorf("Error fetching products: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			ps.Log.Errorf("Error scanning product row: %v", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (ps *ProductsService) GetProductById(id int) (*models.Product, error) {
	ps.Log.Debugf("Fetching product with ID: %d from database", id)

	var product models.Product
	err := ps.DB.QueryRow("SELECT * FROM Products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		ps.Log.Errorf("Error fetching product: %v", err)
		return nil, err
	}

	return &product, nil
}

func (ps *ProductsService) AddProduct(product *models.Product) error {
	ps.Log.Debugf("Creating new product in database, data: %+v", product)

	result, err := ps.DB.Exec("INSERT INTO Products (name, description, price) VALUES (?, ?, ?)", product.Name, product.Description, product.Price)
	if err != nil {
		ps.Log.Errorf("Error creating new product: %v", err)
		return err
	}

	lastInsertID, _ := result.LastInsertId()
	product.ID = int(lastInsertID)

	return nil
}

func (ps *ProductsService) UpdateProduct(id int, product *models.Product) (*models.Product, error) {
	ps.Log.Debugf("Updating product with ID: %d in database, data: %+v", id, product)

	_, err := ps.DB.Exec("UPDATE Products SET name = ?, description = ?, price = ? WHERE id = ?", product.Name, product.Description, product.Price, id)
	if err != nil {
		ps.Log.Errorf("Error updating product: %v", err)
		return nil, err
	}

	// Fetch the updated product from the database
	updatedProduct, err := ps.GetProductById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		ps.Log.Errorf("Error updating product: %v", err)
		return nil, err
	}

	return updatedProduct, nil
}

func (ps *ProductsService) PatchProduct(id int, product *models.Product) (*models.Product, error) {
	ps.Log.Debugf("Patching product with ID: %d in database, data: %+v", id, product)

	// Build the SQL query dynamically based on the provided fields
	var updates []string
	var args []interface{}

	if product.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, product.Name)
	}
	if product.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, product.Description)
	}
	if product.Price != 0 { // Assuming 0 is not a valid price
		updates = append(updates, "price = ?")
		args = append(args, product.Price)
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	updateQuery := fmt.Sprintf("UPDATE Products SET %s WHERE id = ?", strings.Join(updates, ", "))
	args = append(args, id)

	_, err := ps.DB.Exec(updateQuery, args...)
	if err != nil {
		ps.Log.Errorf("Error updating product: %v", err)
		return nil, err
	}

	// Fetch the updated product from the database
	updatedProduct, err := ps.GetProductById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		ps.Log.Errorf("Error updating product: %v", err)
		return nil, err
	}

	return updatedProduct, nil
}

func (ps *ProductsService) DeleteProduct(id int) error {
	ps.Log.Debugf("Deleting product with ID: %d from database", id)

	_, err := ps.DB.Exec("DELETE FROM Products WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrProductNotFound
		}
		ps.Log.Errorf("Error deleting product: %v", err)
		return err
	}

	return nil
}
