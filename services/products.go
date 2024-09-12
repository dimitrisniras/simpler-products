package services

import (
	"database/sql"
	"simpler-products/models"
)

type ProductsServiceInterface interface {
	GetAllProducts() ([]models.Product, error)
	GetProductById()
	AddProduct(product *models.Product) error
	UpdateProduct()
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

func (ps *ProductsService) GetProductById() {
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

func (ps *ProductsService) UpdateProduct() {
}

func (ps *ProductsService) PatchProduct() {
}

func (ps *ProductsService) DeleteProduct() {
}
