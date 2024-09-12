package services

import (
	"database/sql"
	"fmt"
	"simpler-products/models"
)

type ProductsServiceInterface interface {
	GetAllProducts() ([]models.Product, error)
	GetProductById(id int) (*models.Product, error)
	AddProduct(product *models.Product) error
	UpdateProduct(id int, product *models.Product) (*models.Product, error)
	PatchProduct()
	DeleteProduct()
}

type ProductsService struct {
	DB *sql.DB
}

func (ps *ProductsService) GetAllProducts() ([]models.Product, error) {
	rows, err := ps.DB.Query("SELECT * FROM Products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (ps *ProductsService) GetProductById(id int) (*models.Product, error) {
	var product models.Product
	err := ps.DB.QueryRow("SELECT * FROM Products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (ps *ProductsService) AddProduct(product *models.Product) error {
	result, err := ps.DB.Exec("INSERT INTO Products (name, description, price) VALUES (?, ?, ?)", product.Name, product.Description, product.Price)
	if err != nil {
		return err
	}

	lastInsertID, _ := result.LastInsertId()
	product.ID = int(lastInsertID)

	return nil
}

func (ps *ProductsService) UpdateProduct(id int, product *models.Product) (*models.Product, error) {
	_, err := ps.DB.Exec("UPDATE Products SET name = ?, description = ?, price = ? WHERE id = ?", product.Name, product.Description, product.Price, id)
	if err != nil {
		return nil, err
	}

	// Fetch the updated product from the database
	updatedProduct, err := ps.GetProductById(id)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (ps *ProductsService) PatchProduct() {
}

func (ps *ProductsService) DeleteProduct() {
}
